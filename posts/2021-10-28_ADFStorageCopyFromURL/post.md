# Using Azure Storage To Download Internet Files From ADF
*2021-10-28*

## The Problem

ADF has the ability to download files or data from the internet using a web activity or an HTTP linked service.  This works in most cases without issue.  However, recently I came across an issue where an HTTP server-side issue was causing an ADF copy activity using an HTTP linked service to fail. 

In this case, the ADF pipeline's only task was to copy a file from the internet, to Azure Storage, every day. Simple!

<img src="2021-10-28 11_34_00-Drawing1 - Visio Professional.png" width="250"/>

The symptom was the file would only download a *portion* of the file, not the entire file. Confusingly, **No errors were recorded on the Copy Activity.**  It always said successful.

When running ```curl``` on the command line to download the file, we could see that the file was not downloading, and required a retry using the ```Range``` header.  
The error that curl outputted was: 
```
cURL 18 transfer closed with outstanding read data remaining
````
This was obviously something weird with the HTTP server (which was an Apache HTTP server). 

Running ```wget``` did download the file successfully, however it did have to try multiple times and used the HTTP ```Range``` header to continue the download from where it failed.

ADF's copy activity using the HTTP linked service does **not** retry a download if it fails.  It does not have any built in retry capabilities using the ```Range``` HTTP header from what I can tell.  Further, you can not manually specify the HTTP ```Range``` header in the copy activity. 

## The Solution

After trying a few different things, it was apparent that ADF wasn't going to be able to actually complete the HTTP download successfully, in any way.  

The solution was to use a lesser known API action of the Azure Storage Blob REST API: the [Put Blob From URL](https://docs.microsoft.com/en-us/rest/api/storageservices/put-blob-from-url) action.  In the ```Put Blob From URL``` action, you can ask the storage service to save a blob whose source is from an external URL.  
The URL is specified in a header called ```x-ms-copy-source```. When this action is called, a blob is placed in the path specified in the request URI, and asynchronously is copied to that location (optionally can be synchronously).  Thankfully, it appears that the storage service **does** have some retry logic in the HTTP request, and did download the file successfully! 

<img src="great.gif" width="200"/>


So when we jump into the ADF pipeline, we have one *Web Activity* that executes this HTTP request.   The JSON of the pipeline is available [here](pipeline.json).

<img src="2021-10-28 10_34_10-GmAdfTest - Azure Data Factory and 8 more pages - Work - Microsoft​ Edge.png"/>

The details of the Web Activity are shown below:

<img src="2021-10-28 10_38_45-GmAdfTest - Azure Data Factory - Work - Microsoft​ Edge.png"/>

> Note: The Web Activity, when using an HTTP Method of ```PUT```, requires a body in the request.  But, there is no body required in the storage API request.  To get around this, put a dynamic expression of ```@trim('')```.  This gets us past the UI requirement for a body, but an empty body will be sent to the storage API. 

The storage REST API requires that the caller be authenticated.  This can be accomplished by [retrieving an OAuth token](../2021-08-17_GraphAPIAndADF) or by using a *Managed Identity*.  In our case (and in most cases), a Managed Identity should be used. 

### Setting up Managed Identity Permissions

1. An Azure Data Factory has a built in Managed Identity when deployed.  To grant access to the storage account, open the storage account in the portal (this can also be accomplished by using the Azure CLI).  Click ```Access Control (IAM)```.

    <img src="2021-10-28 11_36_26-gmadlshared - Microsoft Azure and 11 more pages - Work - Microsoft​ Edge.png"/>

2. Click Add, then Role Assignment

    <img src="2021-10-28 11_37_04-gmadlshared - Microsoft Azure and 11 more pages - Work - Microsoft​ Edge.png"/>

3. Click Contributor, then Next

    <img src="2021-10-28 11_38_23-Add role assignment - Microsoft Azure and 11 more pages - Work - Microsoft​ Edge.png"/>

4. Select *Managed Identity* for the Assign Access to field.  Search for the Azure Data Factory you want to grant access to.  Then click *Select*.  When done, click *Review + assign*.  

    <img src="2021-10-28 11_39_13-Select managed identities - Microsoft Azure and 11 more pages - Work - Microsoft.png"/>

5. Click *Review + assign*. 

The ADF environment should now have access to call the Storage REST API Put actions. 

> Note: The ADF Managed Identity will have access to all the containers/blobs in that storage account.  You can always be more granular if you wish!

## Summary

Using the ```Put Blob From URL``` action on the Azure Storage REST API allowed us to download a file with retry capabilities successfully, and we were able to call this API method from ADF.  A side benefit from this is, we save the *data movement* costs associated with a standard Copy Activity in ADF.  The only cost associated is the activity run cost, which is negligible.  
