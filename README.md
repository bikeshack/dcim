# Data Center Infrastructure Manager for HPC

This DCIM implementation is an example service based on the Alex Lovell-Troy's [opinionated microservice chassis](https://github.com/alexlovelltroy/chassis) and performs inventory management operations for HPC systems.

## Organization

Code in this example is organized into a data object called a Component in [pkg/components/components.go] and a set of http microservices that manage Create, Retrieve, Update, and Delete operations in [crud.go].  The data is actually stored in a postgres database with an internal library managing that interaction at [internal/postgres/postgres.go].  Instantiating and managing the database schema is handled through the files in the [migrations/] directory.

This code organization relies on an opinionated microservice chassis that provides the command structure.

```

```

## Testing