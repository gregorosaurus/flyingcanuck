# Backing up ADLS Gen2 using AZCopy

## Introduction

Azure Data Lake Storage (ADLS Gen2) is becoming a more integral part of many organization's data strategy.  Many organizations may want to manually manage the backup and retention of this data.  While Azure [Storage Accounts keep three copies of data and have built in redundancy](https://docs.microsoft.com/en-us/azure/storage/common/storage-redundancy) (additionally can copy the data to another data center using Geo Redundancy), companies may want to manually manage or create additional backups. With ADLS Gen2, at the time of writing, not supporting soft deletes, this could be one way organizations can recover from accidental deletion.  Additionally, organizations may want a copy of their data in a separate subscription. 

This post will go through: 
- How to use ```azcopy``` to backup your ADLS Gen2 storage account. 
- How to schedule and run the ```azcopy sync``` from within Azure.

## What about Azure Backup?!

Azure Backup is the primary method of backup for VMs, Storage, and other services within Azure. But unfortunately, at the time of writing, [Azure Backup was **not** available for ADLS Gen2](https://docs.microsoft.com/en-us/azure/backup/blob-backup-support-matrix).  Hopefully in the future this changes. 

<img src="2021-08-31 13_49_07-Support matrix for Azure Blobs backup - Azure Backup _ Microsoft Docs and 5 more.png">

## Prerequisites 

This post assumes the following:
1. A source ADLS Gen2 Storage Account
2. A destination ADLS Gen2 Storage Account (the backup location)
3. AZCopy installed or downloaded

## AZCopy

From Microsoft: AzCopy is a command-line tool that moves data into and out of Azure Storage. We can use this tool to automate the backup of Azure Storage Accounts, including ADLS Gen2.  

Specifically, azcopy includes a command called [sync](https://docs.microsoft.com/en-us/azure/storage/common/storage-ref-azcopy-sync).  This replicates the source location to the destination location.  This command will be the basis of our backup strategy. 

## Authentication

[We can authenticate and authorize azcopy a number of ways](https://docs.microsoft.com/en-us/azure/storage/common/storage-use-azcopy-authorize-azure-active-directory?toc=/azure/storage/blobs/toc.json). For our purposes we can use either SAS Tokens, or Service Principals.  

Service principals are a great option if the storage account you're backing up to is within the same Azure AD Tenancy. If your source and backup storage accounts are in *different* Azure AD tenants, you must use SAS Tokens. 

### Service Principal Setup

If the source and backup storage accounts are within the *same* Azure AD tenant, a single service principal can be setup and used for authentication to each of the storage accounts. 

### SAS Token Setup

If the source and backup storage accounts are **not** within the same Azure AD tenant, you will have to setup and use SAS tokens.  

>These SAS tokens will have expiries, so make sure you note the expiry and update them accordingly.  



## Azure Function

### Timer Trigger

### Timeouts

## Microsoft Azure Storage Data Movement Library