# PCORnet Common Data Model

Official site: http://www.pcornet.org/resource-center/pcornet-common-data-model/

## Additions

- The `REQUIRED` column has been added to denote which set of field values are required for the record to be valid. This set of fields *do not* guarantee uniquess across records in a table.
- The `REF_TABLE` and `REF_FIELD` columns have been added to denote which fields that reference other values. In a relational model, this would be implemented as a foreign key.
