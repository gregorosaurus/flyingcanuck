# Synapse SQL CTE Query Error
*2021-09-09*

## The Error

Recently I received the following error when running a [CTE](https://docs.microsoft.com/en-us/sql/t-sql/queries/with-common-table-expression-transact-sql?view=sql-server-ver15) query in Synapse and was presented with the following error:

```
The query processor could not produce a query plan because the target DML table is not hash partitioned. 
```

## The Fix

The reason for this error is that CTE queries are not supported on non-hash distributed tables.  The table I was running the query against was round-robin distributed.  

My assumption is that there is no uniqueness in the rows, so there is no way to safely delete some records that are affected by the predicate.  But that's just a guess.  

For me, the fix was simply recreate the table as a hash distributed table.  Because this was the early days of the database, this was trivial to do. It was likely an oversight that it was created as a round-robin table, and should have been hashed from the beginning. 