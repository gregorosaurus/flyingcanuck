# Creating a CI/CD Process for Synapse SQL
*2021-10-04*

Synapse and Synapse Pipelines have a deep and thorough integration with git and github. Each change and modification is committed to git which allows for common git patterns such as [gitflow](https://www.atlassian.com/git/tutorials/comparing-workflows/gitflow-workflow).  Further, this git integration makes it relatively straight forward to implement a CI/CD process for promotion into production environments (more on that in future post).  But what happens when you want to include the Synapse SQL Pool database schema with that same CI/CD process?  Sadly, that isn't built into Synapse, and it's something you have to manage yourself. But fret not, there is a way to do this without manually updating the database schemas! 

In this post, we'll review how to migrate changes from one Synapse SQL Pool to another, and how to integrate that into a CI/CD pipeline in github using github actions. 

## Overview


## DACPAC

SQL has long had a mechanism for creating, managing, and updating schemas in a database.  This is a DAC, or [Data-Tier Application](https://docs.microsoft.com/en-us/sql/relational-databases/data-tier-applications/data-tier-applications?view=sql-server-ver15).   For our purposes, we really only need to understand that this is a method of encapsulating all database schema information, including tables, views, and instance objects, for a specific database.  This information is contained within a .dacpac file, which can then be *applied* to other databases to match that schema.  The .dacpac file can be created a number of ways, with a variety of tools.  For our purposes, we will be using the *SQLPackage* tool to both extract and apply a .dacpac. 

## SQL Package

[SQLPackage.exe](https://docs.microsoft.com/en-us/sql/tools/sqlpackage/sqlpackage?view=sql-server-ver15) is a command line utility that automates various tasks against a SQL Server (including Synapse SQL Pools).  This includes extracting schema information from a SQL database, and publishing a schema (DACPAC) against a database.  These are the two tasks we'll use in our CI/CD process. 

To start, let's test to make sure we can run these tasks manually, without an automated CI/CD process. 

### Manual SQLPackage Commands

1. Before we begin, [download SQLPackage](https://docs.microsoft.com/en-us/sql/tools/sqlpackage/sqlpackage-download?view=sql-server-ver15).  
    >If you installed on windows, be sure to at `C:\Program Files\Microsoft SQL Server\150\DAC\bin` into your PATH environment variable, **or** make sure you cd into the directory after opening a command prompt. 
2. Open a command prompt, and go to the directory of where the SQL Package was installed or extracted.
3. Now, the first thing we do is extract a DACPAC from the source database server. In our manual process, we'll use a sql login as the authentication mechanism. 

    ```
    .\SQLPackage.exe /TargetFile:"C:\temp\sql_current_version.dacpac" /Action:Extract /SourceServerName:"[servername].sql.azuresynapse.net,1433" /SourceDatabaseName:"TestDedicatedPool" /SourceUser:adminuser /SourcePassword:"********"
    ```

    <img src="2021-10-04 09_54_55-PowerShell 7 (x64).png"/>
4. After the DACPAC is generated, we can test applying it to another database. This time, we will use the [Publish](https://docs.microsoft.com/en-us/sql/tools/sqlpackage/sqlpackage-publish?view=sql-server-ver15) action.  Again, using the manual method, we'll authenticate using a SQL login.   In our testing, we have an empty SQL Pool.  The expectation is after we run the Publish action, the database will match the schema of the source database. 

    <img src="2021-10-04 10_00_55-SQLQuery1.sql -DestinationTestDedicatedPool.png" width="350"/>

    ```
     .\SqlPackage.exe /TargetFile:"C:\temp\sql_current_version.dacpac" /Action:Publish /TargetServerName:"[servername].sql.azuresynapse.net,1433" /TargetDatabaseName:"TestDedicatedPool" /TargetUser:adminuser /TargetPassword:"********"
    ```

    If everything worked, you should see an output similar to:
    <img src="2021-10-04 10_14_41-post.md - 2021-10-04_SynapseSqlCICD - Visual Studio Code.png"/>

    Now let's check the database using SSMS. 
    <img src="2021-10-04 10_15_52-SQLQuery1.sql - .sql.azuresynapse.net,1433.TestDedicatedPool (gregmar.png" width="300"/>

    Party time!!

    <img src="greatsuccess.gif" width="200"/>

## Automation

Now that we have the process of extracting and publishing a DACPAC using SQLPackage, we can move onto automating this into our CI/CD pipelines. 

For this post, we will use github actions, however it can easily be adapted to Azure DevOps Pipelines.  In fact, there is documentation on how to run [SQLPacakge on Microsoft's website](https://docs.microsoft.com/en-us/sql/tools/sqlpackage/sqlpackage-pipelines?view=sql-server-ver15#additional-sqlpackage-examples). 

Luckily for us, there is a [SQLPackage github action](https://github.com/Azure/run-sqlpackage-action) that saves us some work of downloading and extracting the SQLPacakge zip into the agent environment. We will leverage this as part of our github action. 

> Note: the github action referenced says that only Publish is supported, and doesn't say Extract is supported.  But if you look at the code of the action, you can see that Extract will work fine. 

## Authentication

We have two options for authentication for the CI/CD pipelines.  
1. SQL authentication
2. Service principals. 

SQL authentication is relatively straight forward.  As with the manual steps we took previously, we use a SQL username and password for authentication.  This time however, we must specify this all in a connection string.  For example, the connection string would be something like: 
```
Server=tcp:[servername].sql.azuresynapse.net,1433;Initial Catalog=TestDedicatedPool;Persist Security Info=False;User ID=sqladminuser;Password={your_password};MultipleActiveResultSets=False;Encrypt=True;TrustServerCertificate=False;Connection Timeout=30;
```

For service principals, it's a bit more complex, but likely the better way to go for production CI/CD pipelines. 

### Service Principal Setup

We will create **one** service principal for our CI/CD sql schema updating pipeline.  This service principal will be granted the *db_owner* for the target.

1. Head over to the [Azure Portal](https://portal.azure.com)

2. Open *Azure Active Directory*

3. Click *App Registrations*

    <img src="2021-10-04 10_29_44-Microsoft - Microsoft Azure and 17 more pages - Work - Microsoft​ Edge.png" width="250"/>

4. Click *New registration*
    
    <img src="2021-10-04 10_30_37-Microsoft - Microsoft Azure and 17 more pages - Work - Microsoft​ Edge.png" width="450">

5. Enter in a name for the Service Principal.  After, click Register
    
    <img src="2021-10-04 10_32_30-Register an application - Microsoft Azure and 17 more pages - Work - Microsoft​.png"/>

    **Important**: Write down the Application (client) ID.  This is used later. 

6. Click Certificates & secrets, and then click *New client secret*. 

    <img src="2021-10-04 10_39_35-Add a client secret - Microsoft Azure and 17 more pages - Work - Microsoft​ Edge.png"/>

    **Important**: Write down the secret value.  This will be used later. 

### Github Action

Now we will create the github action. This action will be responsible for 1) extracting the DACPAC from dev/test, and 2) publishing the dacpac to stage/prod/etc. 

> This post assumes that the synapse environment in dev is configured to use git and github.  If not, you can use any github repository to host the workflow.  However, it would be logical to host the workflow in the same repository as the Synapse configured git repository. 

<img src="2021-10-04 10_44_59-Actions · gregorosaurus_SynapseZebra and 4 more pages - Work - Microsoft​ Edge.png" />

