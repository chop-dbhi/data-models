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

## Data Models Service

### Top-level Resources

- Model Specifications - [/models](http://data-models.origins.link/models) - A specification of each data model version is available at a `/models/<data model>/<version>` endpoint (e.g., [/models/omop/v5](http://data-models.origins.link/models/omop/v5)). 

### Content Negotiation

The service supports representing each resource in various formats using simple content negotation. The supported formats are:

- HTML - `text/html`
- Markdown - `text/markdown`
- JSON - `application/json`

The desired format can be requested either by setting the `Accept` header to the corresponding mimetype or by adding a `format` parameter to the URL. For example, below is the OMOP v5 specification resource represented in each format:

- HTML - [/models/omop/v5?format=html](http://data-models.origins.link/models/omop/v5?format=html)
- Markdown - [/models/omop/v5?format=md](http://data-models.origins.link/models/omop/v5?format=md)
- JSON - [/models/omop/v5?format=json](http://data-models.origins.link/models/omop/v5?format=json)

Representations are tailored to the clients that are expected to use the resource, as described below. Note that some resources do not support all formats. The HTML format is the default format provided when neither method of content negotiation are used.

### Model Specification Resources

#### HTML

The HTML format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/v5?format=html)) is intended as a very simple proof of concept for displaying the data model specification in a web client for review by data model and/or data users. As such, it begins with the data model version id and a reference URL, followed by a list of tables (which serves as a linked table of contents). Each table section includes the table description and a list of fields (again, a linked table of contents). For each field, "refers to" information, if it exists, is followed by the description and any schema specifications. A table of mappings and a table of inbound references are also provided, if that information is found. This content represents an aggregation of information about the data model which we think would be useful for data model and/or data users.

#### Markdown

The Markdown format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/v5?format=md)) provides the same information as the HTML format. In fact, the HTML format is derived directly from the Markdown. The specific choices about header levels and organization can be seen at the actual endpoints linked above. This is intended as an API of sorts from which use-case-specific clients can retrieve, process, and display aggregated data model specification information as they wish.

#### JSON

The JSON format (e.g., [OMOP v5](http://data-models.origins.link/models/omop/v5?format=json)), unlike the previously described formats, is intended for technical implementation clients and therefore presents a readily machine-processable and exhaustive representation of the data model specification. The top-level object contains the data model `name`, `version`, and reference `url` as well as an array of `tables`. Each object in the `tables` array contains the table `name` and `description`, an array of `fields`, and the `model` name and model `version`, to unambiguously identify the model to which the table belongs. Each object in the `fields` array contains the field `name`, `description`, `type`, and `required` status (as per governance), as well as the `default` (which defaults to `""`), `length`, `precision`, and `scale` (which all default to `0`). Each field object also contains the `table` name. This format should be useful in dynamically creating all sorts of data model operations, from schema creation to annotation to transformations.
