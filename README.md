## Probum - User Management Service

**Note**: To check each operation for this service, refer to its Swagger documentation ([Swagger Docs](http://127.0.0.1:3000/docs)).

### Service Dependencies:

| Application      | Version  |
|------------------|:--------:|
| `docker`         | ≥ 24.0.7 |
| `docker-compose` | ≥ 1.29.2 |
| `go`             | ≥ 1.21.1 |

### Service Configuration:

Each and every configuration aspect of the service (mainly others endpoints, database settings) can be set in `/config/development.yml`.
The environmental variable values from the `docker-compose` and the configuration file, relative to database settings, must match.

### Running the Service:

Setting up the Postgres database:
```sh
docker-compose up
```
Now the database should be up and running, accessible at port `5432`.
The seeding for the database, happens only once, when the database is empty. This is, the seeding occurs at the first migration.

To re-seed the database, you will need to drop the table for the `users` and re-run the migration sequence.

---
Installing `go`'s dependencies using the `Makefile`:
```sh
make deps
```
---
Running the service:
```sh
make run
```

The default port set on the configuration file, is `3000`. If not changed, the service should be available [here](http://127.0.0.1:3000).

To run the service in a development mode with hot-reload (`CompileDaemon`): 
```sh
make daemon
```
