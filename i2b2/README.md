# i2b2

Data models for the [i2b2](https://www.i2b2.org).

## Additions

- The `Ref Table` and `Ref Field` columns have been added to denote the fields that reference fields on other tables. In a relational model, this would be referred as a foreign key. These columns are for information purposes only since the official i2b2 implementation does not enforce these constraints in the database schema.
