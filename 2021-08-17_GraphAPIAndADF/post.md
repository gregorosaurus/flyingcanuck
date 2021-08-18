# Calling the Microsoft Graph API within Azure Data Factory (ADF)
*August 8th, 2021*

I recently had the request to call the [Microsoft Graph API](https://docs.microsoft.com/en-us/graph/use-the-api) from within Azure Data Factory to store user data within an [Azure Data Lake](https://docs.microsoft.com/en-us/azure/storage/blobs/data-lake-storage-introduction).   There are a couple prerequisites required but overall, is relatively simple process and can easily be added to any ADF pipeline. 

## Prerequisites 
1. An Azure AD Service Principal
2. Delegated permissions to the Graph API for the Azure AD Tenant
3. A location to save Graph API exported data, and a data set to support this saving. 

## Disclaimers
- Key vault is considered out of scope for this tutorial.  Do not store any keys in plain text and use key vault to retrieve the application secret when calling the pipeline created in this keyvault. 
- OData paging will not be covered by this tutorial.  This will be added in a future post probably.  This purpose of this tutorial is to demonstrate the OAuth authentication flow and simple iteration through the data returned by the Graph API. 

## Authentication Flow
The Graph API, like many other web services, uses OAuth 2.0 for authentication.  That means that all requests to the Microsoft Graph API must be authenticated using a valid OAuth mechanism.  The Graph API supports *Delegated Permissions* and *Application Permissions*.  We will use the **Application Permission** token request flow.  

The general flow of the Azure Data Factory pipeline will be as follows:

1. Request an access token using the Azure AD Service Principal
2. Call the Graph API route that has the data you're interested in

## Detailed Steps

### Azure AD Setup
1. Open the [Azure Portal](portal.azure.com)
2. Open the [Azure Active Directory Blade](https://ms.portal.azure.com/#blade/Microsoft_AAD_IAM/ActiveDirectoryMenuBlade/Overview)
3. Click App Registrations 

    <img src="2021-08-17 10_57_26-Automata Solutions - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="200"/>
4. Click New Registration

    <img src="2021-08-17 10_58_16-Automata Solutions - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>
5. Type in an application name and then click Register

    <img src="2021-08-17 10_59_09-Register an application - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

6. Take note of the Application ID (sometimes called the client ID) and the Tenant ID.

    <img src="2021-08-17 10_59_45-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

7. Create a client secret.  This is the secret that will be used to authenticate against the Azure AD OAuth endpoint to request an access token. 

    <img src="2021-08-17 11_00_38-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 11_00_57-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 11_12_49-Add a client secret - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 11_13_23-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge" alt="" width="500"/>

    <img src="2021-08-17 11_13_23-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    >**Important:** You must copy the secret value.  It can only be viewed once. 

8. Now it's time to grant the permissions to the application.  These permissions detail what the app can access, both with delegated permissions and application permissions. 

    <img src="2021-08-17 11_14_11-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    Select the application/permission that you would like to grant.  Likely it's the 
    Graph API, but you can select other permissions if you'd like.
    <img src="2021-08-17 11_14_44-Request API permissions - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 11_15_00-Request API permissions - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    >**Remember:** The application permission set is what we want to grant.  Delegated access is used for interactive logins such as websites or user applications.  ADF will use an application permission as it has no interactive login. 

    Select the specific permission you'd like to grant.  This can be anything that the graph API offers.  For our example, we'll just read all user data.

    <img src="2021-08-17 11_15_32-.png" alt="" width="500"/>

9. Admin authorize the application. 

    <img src="2021-08-17 11_16_20-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 11_16_53-ADF Pipeline - Microsoft Azure - [InPrivate] - Microsoft​ Edge.png" alt="" width="500"/>

### Azure Data Factory Pipeline

Now that we've setup the Azure AD service principal, we can move on to creating the Azure Data Factory pipeline.
> **Note:** A synapse pipeline would work exactly the same.  However, the ARM resource types in the provided JSON pipeline definitions would be slightly different. 

1. Open ADF
2. Create a new pipeline
3. Set up the following pipeline parameters:
    - AppID (String)
    - AppSecret (String)
    - TenantID (String)

    <img src="2021-08-17 11_32_42-GmAdfTest - Azure Data Factory and 19 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>
    
    ```json
    "parameters": {
        "AppID": {
            "type": "string"
        },
        "AppSecret": {
            "type": "string"
        },
        "TenantID": {
            "type": "string"
        }
    },
    ```
4. Create a *Web Activity* that authenticates to Azure AD and requests an access token. Name it OAuthentication. (Json definition below for reference).
    > **Important:** The name of this activity matters.  The name is used as a reference to retrieve the output, specifically the access token. 
    ```json
    {
        "name": "OAuthAuthentication",
        "description": "This step authenticates the service principal to call the graph API",
        "type": "WebActivity",
        "dependsOn": [],
        "policy": {
            "timeout": "0.00:05:00",
            "retry": 1,
            "retryIntervalInSeconds": 30,
            "secureOutput": false,
            "secureInput": false
        },
        "userProperties": [],
        "typeProperties": {
            "url": {
                "value": "@concat('https://login.microsoftonline.com/',pipeline().parameters.TenantID,'/oauth2/v2.0/token')",
                "type": "Expression"
            },
            "method": "POST",
            "headers": {
                "Content-Type": "application/x-www-form-urlencoded"
            },
            "body": {
                "value": "@concat('client_id=',pipeline().parameters.AppID,'&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default&client_secret=',pipeline().parameters.AppSecret,'&grant_type=client_credentials')",
                "type": "Expression"
            }
        }
    }
    ```
5. Under *Settings*, enter in the URL as an expression (Click Add Dynamic Content)

    ```
    @concat('https://login.microsoftonline.com/',pipeline().parameters.TenantID,'/oauth2/v2.0/token')
    ```

6. Add a new header, with the name of ```Content-Type```, and the value as ```application/x-www-form-urlencoded```

7. Set the *Body* to the following expression:
    ```
    @concat('client_id=',pipeline().parameters.AppID,'&scope=https%3A%2F%2Fgraph.microsoft.com%2F.default&client_secret=',pipeline().parameters.AppSecret,'&grant_type=client_credentials')
    ```

8. After all settings are entered, your settings pane should look like the following:

    <img src="2021-08-17 11_40_21-GmAdfTest - Azure Data Factory and 19 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

9. Okay bear with me, this is going to be an adventure into linked services and data sets.  We're going to create a generic HTTP OAuth linked service, and then a Graph API data set.  Using that data set, we'll invoke a copy activity to actually copy the data to the data lake. 

    So let's start with the linked service. 

    <img src="2021-08-17 12_43_10-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

10. Create the Generic HTTP Linked service.

    <img src="2021-08-17 12_45_20-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 12_44_28-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

11. Create the Generic Graph API Data set

    <img src="2021-08-17 12_46_48-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 12_48_23-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

    <img src="2021-08-17 12_48_41-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

12. Set the BaseURL parameter to ```https://graph.microsoft.com```
    
13. Set the Relative URL to ```@dataset().GraphAPIRoute```
    
    <img src="2021-08-17 12_48_01-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

14. Create a parameter:

    <img src="2021-08-17 12_51_26-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

15. Okay. we should be done with the linked service and data set. 
Here is the data set JSON, for reference. 
```json
{
    "name": "HTTP_GraphAPI_Data",
    "properties": {
        "linkedServiceName": {
            "referenceName": "HTTP_Generic_LinkedService",
            "type": "LinkedServiceReference",
            "parameters": {
                "BaseURL": {
                    "value": "https://graph.microsoft.com",
                    "type": "Expression"
                }
            }
        },
        "parameters": {
            "GraphAPIRoute": {
                "type": "string",
                "defaultValue": "/v1.0/users"
            }
        },
        "folder": {
            "name": "BuiltIn"
        },
        "annotations": [],
        "type": "Json",
        "typeProperties": {
            "location": {
                "type": "HttpServerLocation",
                "relativeUrl": {
                    "value": "@dataset().GraphAPIRoute",
                    "type": "Expression"
                }
            }
        },
        "schema": {}
    },
    "type": "Microsoft.DataFactory/factories/datasets"
}
```

16. Go back to our pipeline, and create a *copy activity*. This copy activity will be the thing doing the actual work to get our graph response to the data lake. 

    <img src="2021-08-17 12_54_41-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

17. Set the source to our HTTP Graph API generic data set. Set the GraphAPIRoute parameter to the route you'd like to call. In our case, it'll be ```/v1.0/users```. Set the additional headers to the following expression:
    ```
    Authorization: Bearer @{activity('OAuthAuthentication').output.access_token}
    ```

    <img src="2021-08-17 12_56_23-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

18. Set the sink to the destination you'd like to save the file. In our case, we have a generic CSV dataset that we will use.  (You can use anything though).  I've included the JSON for this dataset for reference:
    ```json
    {
        "name": "CSV_File",
        "properties": {
            "linkedServiceName": {
                "referenceName": "AzureDataLakeShared",
                "type": "LinkedServiceReference"
            },
            "parameters": {
                "Container": {
                    "type": "string"
                },
                "Directory": {
                    "type": "string"
                },
                "File": {
                    "type": "string"
                }
            },
            "folder": {
                "name": "BuiltIn"
            },
            "annotations": [],
            "type": "DelimitedText",
            "typeProperties": {
                "location": {
                    "type": "AzureBlobFSLocation",
                    "fileName": {
                        "value": "@dataset().File",
                        "type": "Expression"
                    },
                    "folderPath": {
                        "value": "@dataset().Directory",
                        "type": "Expression"
                    },
                    "fileSystem": {
                        "value": "@dataset().Container",
                        "type": "Expression"
                    }
                },
                "columnDelimiter": ",",
                "escapeChar": "\\",
                "firstRowAsHeader": true,
                "quoteChar": "\""
            },
            "schema": []
        },
        "type": "Microsoft.DataFactory/factories/datasets"
    }
    ```

    <img src="2021-08-17 12_57_54-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="750"/>

19. Set the mapping of the copy.  In our case, we must manually specify the mapping.  Depending on your destination dataset type, you may not have to. 

    <img src="2021-08-17 12_59_38-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>   
    
    > Using the [Microsoft Graph Explorer](https://developer.microsoft.com/en-us/graph/graph-explorer) can help determine which fields you'd like to save and map. 

## Testing
Okay! We have a pipeline that looks generally like this:

<img src="2021-08-17 12_54_41-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

Now we need to test it.  Hit that big ol' debug button and enter in the values for the service principal.  
> **Important:** In production, create a pipeline that executes the pipeline that we created with the parameters set, and use key vault to retrieve the application secret. 

<img src="2021-08-17 13_03_08-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

Wait for it to complete, and let's check out the data lake location! 

<img src="2021-08-17 13_04_00-GmAdfTest - Azure Data Factory and 18 more pages - Work - Microsoft​ Edge.png" alt="" width="500"/>

Off to the data lake!
There it is!

<img src="2021-08-17 13_04_47-development - Microsoft Azure - Work - Microsoft​ Edge.png" alt="" width="500"/>

A quick preview:

<img src="2021-08-17 13_05_29-raw_graph_users.csv - Microsoft Azure - Work - Microsoft​ Edge.png" alt="" width="500"/>

And there we have it! 
> **Important:** Make sure to publish your changes so you don't lose them!

## Summary

So what did we do?  We created a pipeline that 1) Authenticated with Azure AD using an Azure AD Service Principal, 2) created a generic Graph API dataset that used an HTTP linked service to 3) download graph user data and 4) save it into an ADLS gen2 data lake.  

Now you can use this data in any transformations, or other data warehouse tasks that you'd like. Using the Graph API from within ADF allows you to centrally manage these web requests like any other data service within your orgnization, benifiting from the capabilities like monitoring and scale that comes with using ADF.  

