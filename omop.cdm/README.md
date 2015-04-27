# OMOP Common Data Model

Official site: http://omop.org/CDM

## Additions

- The `Ref Table` and `Ref Field` columns have been added to denote which fields that reference other values. In a relational model, this would be implemented as a foreign key.

## Notes

- Cohort `subject_id` can point to a Person, Provider, or Visit Occurrence record. It is currently represented in the data model.
