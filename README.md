# Data Center Infrastructure Manager for HPC

This DCIM implementation is an example service based on the Alex Lovell-Troy's [opinionated microservice chassis](https://github.com/alexlovelltroy/chassis) and performs inventory management operations for HPC systems.

## Organization

Code in this example is organized to illustrate separation of concerns as an example of building a testable microservice based on an opinionated chassis.

### The Chassis
The first concept is the [opinionated chassis](https://github.com/alexlovelltroy/chassis) that arranges [gin](https://gin-gonic.com/), [sqlx](http://jmoiron.github.io/sqlx/), [migrate v4](https://github.com/golang-migrate/migrate), [cobra](https://cobra.dev/) in a structure that makes it easy to obey the standards and possible to go beyond.

## Separation of Concerns
The second idea central to the code in this repo is the separation between the data object itself, the database representations of that object, and the CRUD APIs that manage the object.

[pkg/components/compontents.go](/pkg/components/components.go) defines our data object with struct tags to standardize json and sql representations.  In addition, the Component object has methods for validation and for conversion across formats.

Because the Component object is self-contained, it can be imported by other go code for easy serialization/deserialization and standard validation without calling upstream services.

[internal/postgres/postgres.go](/internal/postgres/postgres.go) defines a simple SQL representation of the Component object within a Postgres database.  This is stored in `/internal` rather than `/pkg` because it is an implementation detail of the service and can be changed without impact on the data object itself or the api used by the CRUD operations.  Schema management of the postgres tables is handled through the chassis with a predefined `migrate` command.  Example migrations which set up the database tables are in the [migrations](/migrations/) directory.

Finally, the CRUD handlers themselves are part of the `main` package of the service.  But all business logic is isolated from runtime considerations within the [crud.go](/crud.go) file which defines and then implements a chassis-based microservice.  Critically, each handler is independently testable, even without database access as shown in [crud_test.go](/crud_test.go).

## Runtime Concerns

The binary generated from the code in this repository is suitable for running independently on the commandline.  It accepts flags for configuration and is prepared for use within a container.  The included [Dockerfile](/Dockerfile) builds a minimal container image with [tini](https://github.com/krallin/tini) as the entrypoint.

### Configuration

This software rejects configuration files in favor of flags and environment variables which are more suitable for containerized applications.

