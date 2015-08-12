# Data Models

Data models and vocabularies in the biomedical space.

## Persistent CSV Format

Data model descriptions are stored persistently in this repository in CSV format for portability and human readability. Each data model has its own directory with versions of the model in subdirectories. Each version directory has a `datamodel.json` file that holds metadata about the datamodel and version, so as not to rely on directory structure for interpretability. In fact, this file and a collection of CSV files with the below described header signatures is enough to signal that a data model definition exists. However, the organization and naming conventions presented below have been useful in our initial data model definitions.

Each data model version should have at least `definitions` and `schema` directories and, optionally, a `constraints` directory and `indexes.csv` and `references.csv` files.

The `definitions` directory (e.g., [omop/v5/definitions](omop/v5/definitions)) holds basic information about the data model that would be of primary interest to a data user. There is a `tables.csv` file (e.g., [omop/v5/definitions/tables.csv](omop/v5/definitions/tables.csv)), which lists `name` and `description` for each table, as well as a CSV file for each table (e.g., [omop/v5/definitions/person.csv](omop/v5/definitions/person.csv)), which lists `name` and `description` for each field, whether the field is `required` (a governance, not schema, attribute), and optionally a `ref_table` and `ref_field` combination to which the field refers (typically manifested as a foreign key relationship).

The `schema` directory holds detailed information that might be used to instantiate the data model in a database or other physical storage medium. There is a CSV file for each table (e.g., [omop/v5/schema/person.csv](omop/v5/schema/person.csv)) that lists `type`, `length`, `precision`, `scale`, and `default` attributes (all optional except `type`) for each field, which is identified by `model`, `version`, `table` name, and `field` name attributes.

The `constraints` directory (e.g., [omop/v5/constraints](omop/v5/constraints)), if present, can hold any number of CSV files which list data level constraints that should be applied to any physical representation of the data model. These files (e.g., [omop/v5/constraints/not_nulls.csv](omop/v5/constraints/not_nulls.csv)) contain a `type`, an optional `name`, and the target `table` and `field` for each constraint.

The `indexes.csv` file (e.g., [omop/v5/indexes.csv](omop/v5/indexes.csv)), if present, lists indexes that should be built on a physical representation of the data model, with `name`, whether the index should be `unique`, target `table` and `field`, and `order` attributes for each index.

The `references.csv` file (e.g., [omop/v5/references.csv](omop/v5/references.csv)), if present, lists references (usually foreign keys) which should be enforced on the data model. Each reference is listed with the source `table` and `field`, the target `table` and `field`, and an optional `name`.

Each data model root directory may have a `renamings.csv` file (e.g., [omop/renamings.csv](omop/renamings.csv)) that maps fields which have been renamed across versions by providing a source data model `version`, `table`, and `field` and a target `version`, `table`, and `field`.

The top-level `mappings` directory holds a series of CSV files which list field level mappings between data models. The files (e.g., [mappings/pedsnet_v2_omop_v5.csv](mappings/pedsnet_v2_omop_v5.csv)) contain a `target_model`, `target_version`, `target_table`, and `target_field` as well as a `source_model`, `source_version`, `source_table`, and `source_field` along with a free text `comment` for each mapping.

## CSV Tools

#### Python

The [`csv`](https://docs.python.org/2/library/csv.html) can be used in the standard library.

```python
import csv

# Writes all records to a file given a filename, a list of string representing
# the header, and a list of rows containing the data.
def write_records(filename, header, rows):
    with open('person.csv', 'w+') as f:
        w = csv.writer(f)

        w.writerow(header)

        for row in rows:
            w.writerow(row)
```

#### PostgreSQL

PostgreSQL provides valid CSV output using the [`COPY`](http://www.postgresql.org/docs/9.2/static/sql-copy.html) statement. The output can be to an file using an absolute file name or to STDOUT.

Absolute path.

```sql
COPY ( ... )
    TO '/path/to/person.csv'
    WITH (
        FORMAT csv,
        DELIMITER ',',
        NULL '',
        HEADER true,
        ENCODING 'utf-8'
    )
```

To STDOUT.

```sql
COPY ( ... )
    TO STDOUT
    WITH (
        FORMAT csv,
        DELIMITER ',',
        NULL '',
        HEADER true,
        ENCODING 'utf-8'
    )
```


##### Java

The [`opencsv`](http://opencsv.sourceforge.net/) is a popular package for reading and writing CSV files.

For loop with `rows` as a Collection or Array.

```java
CSVWriter writer = new CSVWriter(new FileWriter(fileName),
                                 CSVWriter.DEFAULT_SEPARATOR,
                                 CSVWriter.NO_QUOTE_CHARACTER);

writer.writeNext(header)

for (int row : rows) {
    writer.writeNext(row);
}

writer.close();
```

If `rows` is a `java.sql.ResultSet`, use `writeAll` directly.

```java
CSVWriter writer = new CSVWriter(new FileWriter(fileName),
                                 CSVWriter.DEFAULT_SEPARATOR,
                                 CSVWriter.NO_QUOTE_CHARACTER);

// Pass the result set and derive the header from the result set
// (assuming it is valid with the spec).
writer.writeAll(rows, true);

writer.close();
```

##### Oracle

Oracle experts should feel free to chime in, but a very promising option is Oracle's new SQLcl command-line tool, available on an early-adopter basis as part of the [SQL Developer](http://www.oracle.com/technetwork/developer-tools/sql-developer/downloads/index.html) family.  SQLcl is being touted as a modern replacement for SQL*Plus.

Sample usage:

```
set sqlformat csv
spool footable.csv
select * from footable;
spool off
```

Another option is to use the [SQL Developer](http://www.oracle.com/technetwork/developer-tools/sql-developer/downloads/index.html) GUI itself, which, although convenient, is not amenable to automation, as SQLcl is.

SQL Developer (and probably SQLcl) export CSV using the following conventions: all text fields are wrapped in quotes (even NULL values, because NULL and empty string are treated the same in Oracle), and no numeric fields are wrapped in quotes. Quotes within fields are escaped via doubling. Newlines within fields are included in the output.

SQL Developer usage:

* On a Data tab (or a table name in the Connections panel), right-click and choose Export
* Change format to csv
* Change line terminator to Unix - other formatting and encoding defaults are fine
