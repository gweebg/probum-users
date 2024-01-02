## Probum - User Management Service

Probum user management service, allows for user creation, deletes, updates. This service does not handle user authentication or
authorization, refer to the Authentication service for that. Password hashes are also not stored in this service database.

**Note**: To check each operation for this service, refer to its Swagger documentation ([Swagger Docs](http://127.0.0.1:3000/swagger/index.html)).

### Service Dependencies:

To compile and run the service, you must have present the following dependencies:

| Application      | Version  |
|------------------|:--------:|
| `docker`         | ≥ 24.0.7 |
| `docker-compose` | ≥ 1.29.2 |
| `go`             | ≥ 1.21.1 |

### Service Configuration:

Before starting the service, you must configure the application via the configuration file located at `/config/development.yml`. This 
allows for database secret definition and other microservice related settings.

The environmental variable values from the `docker-compose` and the configuration file, relative to database settings, must match.

Here, follows an example configuration file: 
```yaml

# Database configuration
db:
  host: "host"
  user: "user"
  password: "password"
  dbname: "userdb"
  port: 5432
  tz: "Europe/Lisbon"

app:
  listen: ":3000"
  jwt-secret: "your-jwt-secret"

```

### Running the Service:

First of all, we need to set up the Postgres database with the command:
```sh
docker-compose up
```
Now the database should be up and running, accessible at port `5432`.
The seeding for the database, happens only once, when the database is empty. This is, the seeding occurs at the first migration.

To re-seed the database, you will need to drop the table for the `users` and re-run the migration sequence.

---
Next, we install `go`'s dependencies:
```sh
go mod tidy
```

---
Before running the service we should first generate the documentation:
```sh
make run
```

The default port set on the configuration file, is `:3000`. If not changed, the service should be available [here](http://127.0.0.1:3000).

(Optional) To run the service in a development mode with hot-reload (`CompileDaemon`): 
```sh
make daemon
```

**Note:** Both commands above generate the `swagger` documentation for the API.
