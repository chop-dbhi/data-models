# Data Models Service

### Top-level Resources

- Models - [/models](/models)

### Content negotiation

The service supports representing each resource in various formats using simple content negotation. The supported formats are:

- HTML - `text/html`
- Markdown - `text/markdown`
- JSON - `application/json`

The desired format can be requested either by setting the `Accept` header to the corresponding mimetype or by adding a `format` parameter to the URL. For example, below is the OMOP v5 resource represented in each format:

- [HTML](/models/omop/v5?format=html)
- [Markdown](/models/omop/v5?format=md)
- [JSON](/models/omop/v5?format=json)

Representations are tailored to the clients that are expected to use the resoure. For example, the JSON representation will typically consist of a lot more information since a consumer of the resource typically uses it to drive a separate client. Note that some resources do not support all formats.
