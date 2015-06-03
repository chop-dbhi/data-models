# Data Models Service

### Top-level Resources

- Model Specifications - [/models](/models) - A specification of each data model version is available at a `/models/<data model>/<version>` endpoint (e.g., [/models/omop/v5](/models/omop/v5)). 

### Content negotiation

The service supports representing each resource in various formats using simple content negotation. The supported formats are:

- HTML - `text/html`
- Markdown - `text/markdown`
- JSON - `application/json`

The desired format can be requested either by setting the `Accept` header to the corresponding mimetype or by adding a `format` parameter to the URL. For example, below is the OMOP v5 resource represented in each format:

- HTML - [/models/omop/v5?format=html](/models/omop/v5?format=html)
- Markdown - [/models/omop/v5?format=md](/models/omop/v5?format=md)
- JSON - [/models/omop/v5?format=json](/models/omop/v5?format=json)

Representations are tailored to the clients that are expected to use the resource, as described below. Note that some resources do not support all formats. The HTML format is the default format provided when neither method of content negotiation are used.

### Model Specification Resources

#### HTML

The HTML format (e.g., [OMOP v5](/models/omop/v5?format=html)) is intended as a very simple proof of concept for displaying the data model specification in a web client for review by data model and/or data users. As such, it begins with the data model version id and a reference URL, followed by a list of tables (which serves as a linked table of contents). Each table section includes the table description and a list of fields (again, a linked table of contents). For each field, "refers to" information, if it exists, is followed by the description and any schema specifications. A table of mappings and a table of inbound references are also provided, if that information is found. This content represents an aggregation of information about the data model which we think would be useful for data model and/or data users.

#### Markdown

The Markdown format (e.g., [OMOP v5](/models/omop/v5?format=md)) provides the same information as the HTML format. In fact, the HTML format is derived directly from the Markdown. The specific choices about header levels and organization can be seen at the actual endpoints linked above. This is intended as an API of sorts from which use-case-specific clients can retrieve, process, and display aggregated data model specification information as they wish.

#### JSON

The JSON format (e.g., [OMOP v5](/models/omop/v5?format=json)), unlike the previously described formats, is intended for technical implementation clients and therefore presents a readily machine-processable and exhaustive representation of the data model specification. The top-level object contains the data model `name`, `version`, and reference `url` as well as an array of `tables`. Each object in the `tables` array contains the table `name` and `description`, an array of `fields`, and the `model` name and model `version`, to unambiguously identify the model to which the table belongs. Each object in the `fields` array contains the field `name`, `description`, `type`, and `required` status (governance-level), as well as the `default` (which defaults to `""`), `length`, `precision`, and `scale` (which all default to `0`). Each field object also contains the `table` name. This format should be useful in dynamically creating all sorts of data model operations, from schema creation to annotation to transformations.
