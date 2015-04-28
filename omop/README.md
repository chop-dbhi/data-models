# OMOP Common Data Model

Official site: http://omop.org/CDM

## Additions

- The `Ref Table` and `Ref Field` columns have been added to denote which fields that reference other values. In a relational model, this would be implemented as a foreign key.

## Notes

### Version 4

- Cohort `subject_id` can point to a Person, Provider, or Visit Occurrence record. It is currently represented in the data model.


### Version 5

- observation table has been removed?
- fact_relationship table `fact_id_1` and `fact_id_2`?
