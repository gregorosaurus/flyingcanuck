# Manually Specifying Synapse Mapping Types for JSON Datasets
*2021-03-15*

## Situation
A customer had created a pipeline that copied data from an REST API to an Azure SQL Database.  They wanted to create the table using *auto table create* feature in the Synapse Copy Activity.  They did specify a mapping for the data to flatten it, but synapse was showing an error when running the pipeline:

```json
{ "errorCode": "2200", "message": "ErrorCode=UserErrorSchemaMappingCannotInferSinkColumnType,'Type=Microsoft.DataTransfer.Common.Shared.HybridDeliveryException,Message=Data type of column 'comment' can't be inferred from 1st row of data, please specify its data type in mappings of copy activity or structure of DataSet.,Source=Microsoft.DataTransfer.Common,'", "failureType": "UserError", "target": "RedactedDataset JSON to DB", "details": [] }
```

The error was occurring because the first json object in the array contained null fields, and when Synapse was trying to create the table in the database, it couldn't determine the data type of the column. 

## Solution
In the UI, it isn't possible to set the data types (hopefully this changes in the future).  However, you *can* set the data type manually in the copy activity json. In each data you're bringing over to the SQL database, you must add the `type` json property in the source and sink, as shown below:

<img src="2022-03-15 11_55_03-RE_jsonTypeSpecification.png" />

If we go look at the copy activity once more, and turn on the advanced editor, we can see that the String data type is specified. 

**IMG**

The following types are available:
* String
* Integer
* Boolean

