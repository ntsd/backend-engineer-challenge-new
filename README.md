# Backend Engineer Coding Challenge

## API Architecture

```mermaid
C4Component
    title API Architecture

    Person(client, "User")

    Boundary(process, "Go Process") {
        Component(api, "API", "Go", "Provides all the API requests.")

        ComponentQueue(scanner, "Scanner", "Go Routine", "A worker pool for repository scanner.")
    }

    ComponentDb(db, "Database", "Postgres", "Store everything here")

    Rel_D(client, api, "HTTP request")
    Rel_R(api, scanner, "Create", "Go routine")
    Rel_D(scanner, db, "Uses")
    Rel_D(api, db, "Uses")
```

## Requirements

- [Go 1.8 or above](https://go.dev/doc/install)

## Libraries

This package uses the following libraries to help with the development.

- github.com/gin-gonic/gin - a HTTP web framework
- github.com/golang/mock - a mocking framework
- gorm.io/gorm - a ORM library
- gorm.io/driver/postgres - a Postgres driver for Gorm
- github.com/go-gormigrate/gormigrate/v2 - database migration for GORM
- github.com/go-git/go-git - Git implementation for Go
- github.com/sirupsen/logrus - a structured logger for Go

## Project structure

The project structure is following [peoject-layout](https://github.com/golang-standards/project-layout) non official standard.

```tree
.
├── cmd:                            - main command files
├── deployments:                    - deployment files
├── internal
│   ├── config:                     - environment variables config
│   ├── storage:                    - database connector and ORM query
│   ├── handler:                    - HTTP handlers
│   ├── model:                      - models of the database
│   ├── scanner:                    - repository scanner worker
```

## Database design

```mermaid
erDiagram
    Repository {
        string id pk
        string name "repository name"
        string url "repository URL"
        time created_at "created time"
        time updated_at "latest updated time"
        time deleted_at "soft deleted time"
    }
    Scan {
        string id pk
        string repository_id fk "repository ID"
        string repository_url "repository URL will not change after origin URL is changed"
        scan_status status "Queued|In Progress|Success|Failure"
        json findings "scan findings json object"
        time started_at "scan started time"
        time finished_at "scan finished time"
        time created_at "scan created time"
        time updated_at "latest updated time"
        time deleted_at "soft deleted time"
    }
    Repository ||--o{ Scan : ""
```

## Environment Variables

| Key          | Description                             | Example                                       |
| ------------ | --------------------------------------- | --------------------------------------------- |
| APP_PORT     | port number, default: `8080`            | 8080                                          |
| POSTGRES_URL | **Required** PostgreSQL URL             | postgresql://user:pass@localhost:5432/db_name |
| SCAN_WORKERS | number of scanner workers, default: `2` | 2                                             |

## Run (Docker)

This package has the Dockerfile and the Docker compose provide to make it easier to run.

```sh
docker-compose build # or `make docker-build`
docker-compose up # or `make docker-up`
```

## Run local

```sh
make run
# or
go run cmd/run/main.go
```

## Run unit test

This package using gomock for the mocking framework.

```sh
make test
# or
go test ./...
```

## Generate mock

To generate mock require [gomock](https://github.com/golang/mock) mockgen

```sh
go install github.com/golang/mock/mockgen@v1.6.0
```

cmd to generate mock

```sh
make mockgen
```

## API Reference

### Entity

the json response entity will following type

#### Success Response

success response will response in json format

| field | type   | description                                                                                  |
| ----- | ------ | -------------------------------------------------------------------------------------------- |
| data  | object | response object can be [Repository](#Repository), [[Repository](#Repository)], [Scan](#Scan) |

#### Error Response

http response status code other than `200` will be considor as an error

| field   | type | description          |
| ------- | ---- | -------------------- |
| message | text | message of the error |

possible http status and error message

| HTTP status | message              | description                      |
| ----------- | -------------------- | -------------------------------- |
| 400         | bad request          | something wrong with the request |
| 404         | not found            | data not found                   |
| 422         | data is duplicated   | some of field is unique          |
| 500         | something went wrong | internal server error            |

#### Repository

Git repository object. the URL need to be unique

*Note: only the Get Repository endpoint will allow to query the scan results for performance purpose

| field     | type            | description                                                                  |
| --------- | --------------- | ---------------------------------------------------------------------------- |
| id        | text            | id of the repository                                                         |
| name      | text            | name of the repository                                                       |
| url       | text            | url of the repository                                                        |
| createdAt | time            | the created time of the repository                                           |
| updatedAt | time            | the latest updated time of the repository                                    |
| scans     | [[Scan](#Scan)] | scan list of the repository, you can only get by get one repository endpoint |

#### Scan

Git scan result object will create every time you scan the repo

| field         | type                  | description                                                    |
| ------------- | --------------------- | -------------------------------------------------------------- |
| id            | string                | id of the scan                                                 |
| repositoryID  | string                | repositoryID will reference to Repository                      |
| repositoryURL | string                | repositoryURL will not change after the repository url changed |
| status        | string                | scan status (Queued, In Progress, Success, Failure)            |
| findings      | [[Finding](#Finding)] | finding object after scan found something                      |
| startedAt     | string                | the started time of the scan                                   |
| finishedAt    | string                | the finished time of the scan                                  |
| createdAt     | string                | the created time of the scan                                   |
| updatedAt     | string                | the latest updated time of the scan                            |

#### Finding

list of scan when found something

| field       | type   | description            |
| ----------- | ------ | ---------------------- |
| path        | string | code path              |
| line        | number | line of code           |
| description | string | detail for the finding |

### Endpoints

#### Create Repository

an endpoint to create Git repository

`POST /repositories`

Request body

| field | type | description            |
| ----- | ---- | ---------------------- |
| name  | text | name of the repository |
| url   | text | url of the repository  |


Response: [Repository](#Repository)

Example

```sh
curl http://localhost:8080/repositories \
    --header 'Content-Type: application/json' \
    --data '{
        "name": "scan-test",
        "url": "https://github.com/ntsd/scan-test"
    }'
```

#### List Repositories

an endpoint to list all repositories

`GET /repositories`

Response: array of [Repository](#Repository)

Example

```sh
curl http://localhost:8080/repositories
```

#### Get Repository

get one repository, this will allow to get the scan results

`GET /repositories/:repositoryID`

Response: [Repository](#Repository) include list of [Scan](#Scan)

Example

```sh
curl http://localhost:8080/repositories/f775e48a-1515-411c-b004-994ec66e1111
```

#### Update Repository

patch will allow to update the repository by fields, no need to add the full repository object

`PATCH /repositories/:repositoryID`

Request body

| field | type | description            |
| ----- | ---- | ---------------------- |
| name  | text | name of the repository |
| url   | text | url of the repository  |

Response: an updated [Repository](#Repository)

Example

```sh
curl --request PATCH http://localhost:8080/repositories/f775e48a-1515-411c-b004-994ec66e1111 \
    --header 'Content-Type: application/json' \
    --data '{
        "name": "scan-test-2"
    }'
```

#### Delete Repository

delete repository

`DELETE /repositories/:repositoryID`

Example

```sh
curl --request DELETE 'http://localhost:8080/repositories/f775e48a-1515-411c-b004-994ec66e1111'
```

#### Scan Repository

scan repository

`POST /repositories/:repositoryID/scan`

Response: [Scan](#Scan)

Example

```sh
curl --request POST 'http://localhost:8080/repositories/f775e48a-1515-411c-b004-994ec66e1111/scan'
```

## Improvements

an improvement list that can't finish in an amount of time.

- separate scanner worker to another instance for a better scaling.
- clone repository to temporary directory instead of re clone every scan, and git reset to remote origin.
- better queue system, maybe implement a 3rd party tool.
