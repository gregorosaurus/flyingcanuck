# Authorizing Managed Identities in Azure Data Explorer
*2022-04-19*

Azure Data Explorer (ADX) has Azure AD authentication and authorization built into the product.  Users can belong to one or many roles defined by an Azure AD user account, an Azure AD group, or an [Azure AD service principal](https://docs.microsoft.com/en-us/azure/data-explorer/provision-azure-ad-app).  

But what happens if you want to use a [Managed Identity](https://docs.microsoft.com/en-us/azure/active-directory/managed-identities-azure-resources/overview) that's used by an Azure Data Factory or Synapse?  Well not to fear, it's extremely straight forward and very similar to authorizing a service principal as part of an ADX role, except you don't have to create the service principal.  

To get started, we need two pieces of information: 

1. The MSI Azure AD Object ID (a GUID representing the application in Azure AD). 
2. The Tenant ID (another GUID representing the Azure AD tenant). 

Let's use an Azure Data Factory for our example.  

1. Locate the ADF resource in the Azure portal. 
<img src="2022-04-19 17_16_02-Data factories - Microsoft Azure - Work - Microsoft​ Edge.png" />

2. Click Properties under the Settings section.
<img src="2022-04-19 17_17_28-GmAdfTest - Microsoft Azure - Work - Microsoft​ Edge.png" />

3. Note the Managed Identity Object ID and the Managed Identity Tenant ID. 
<img src="2022-04-19 17_18_27-GmAdfTest - Microsoft Azure - Work - Microsoft​ Edge.png" />

4. Open your [ADX cluster](https://dataexplorer.azure.com) and open a new query. 
<img src="2022-04-19 17_21_12-adxgmsharedprd.westus3.shared _ Azure Data Explorer - Work - Microsoft​ Edge.png" />


Congratulations, your MSI should now have access to your ADX cluster. 
> In this example we are using the `viewers` role.  Be sure to check out the [roles available in ADX](https://docs.microsoft.com/en-us/azure/data-explorer/kusto/management/security-roles#managing-database-security-roles).

<img src="missionaccomplished.jpg" width="400" />