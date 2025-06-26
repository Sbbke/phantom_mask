# System Arachitecture

## Project Structure ##

```bash
/data
├── pharmacies.json
├── user.json
/docker
├── .env                       // define environment virables
├── docker-compose-db.yaml     // init database container
├── docker-compose-init.yaml   // data process and migration
├── docker-compose.yaml        // run the service
├── servers.json
/Phantom_BE
├── Dockerfile                 // containerize phantom_be service
├── go.mod
├── main.go                    // call cmd.Execute()
└── cmd/
    ├── root.go                // define Cobra entries
└── global/                    // define global vriables, const, and struct
    ├── common.go              // global constant value
    ├── global.go              // global struct and vriables
└── app/                       // backend application
    ├── app.go                 // main entry point
    └── api/          
        ├── pharmacy_dto.go    // define dto of pharmacy api
    └── controllers/           // define handlers for router path
        ├── hello.go  
        ├── pharmacy_controller.go  // define api entries for pharmacy
        ├── pharmacy_helper.go      // define helper funciton for pharmacy api
    └── initial/              
        ├── initial.go         // data initialization 
        ├── etl_helper.go      // define function for data preprocessing
        ├── etl_helper_test.go // test etl function
    └── middleware/            // define custom middleware handler
        ├── common.go          // add common rules in middleware
        ├── recovery.go        // add handler for panic recovery, database error and timeout
        ├── rateLimie.go       // (empty) Define ratelimit rules
    └── models/              
        ├── database.go        // control over database, from connection to migration
    └── routes/
        ├── router.go          // define router rules
    └── validation/
        ├── validator.go       // define validator for pharmacy api endpoints
└── test/    
    └── integration/           // containe test file
```
## Components ##
### API server ###
+ Uses Gin for HTTP routing

+ Implement custom middleware handler for:
    - Common header
    - Panic recovery
    - Database error
    - Timeout
+ Uses Cobra CLI for database migration & ETL commands
+ GORM for ORM
+ Exposes REST endpoints (e.g., /api/v1/pharmacies, /api/v1/hello)
+ Validates incoming requests via go-playground/validator

### Postgrsql Database ###

+ Pharmacy Database: Stores pharmacies, masks, users, purchase histories

+ User Database: Stores login/authentication data (optional use)

### ETL Logic ###
+ Parses raw JSON data from /data/pharmacies.json and /data/user.json

+ Normalizes fields such as openingHours into structured formats

+ Performed in Go (initial/etl_helper.go)

### Data Migration ###
+ Insert parsed data into PostgreSQL tables

+ Triggered via CLI commands:       
    - > migrateSchema
    - > initUsers
    - > initPharmacies

### Docker Services ###
+ postgres-user: User DB container

+ postgres-pharmacy: Pharmacy DB container

+ phantom-be: Go backend container

+ pgadmin: DB management UI

To know how to run in docker environment, please click [here](./deployment.md)



## ETL and Migration ##

1. Separating preprocessing and loading (into database) may be better if:

+ Dataset grows significantly, causing performance issues.
+ Need to reuse preprocessed data elsewhere (e.g., for analytics).
+ Want to isolate transformation errors from database errors.

2. Models.MigrateSchema uses GORM’s AutoMigrate. For production, switch to golang-migrate for versioned, idempotent migrations.


3.To further simplify, could combine initUserSchema and initPharmaciesData into a single Cobra command (e.g., initData).

## Testing ##

For testing pharmacy cotrollers' business logic, I use integration tests instead unit test. 

Since I uses gorm directly inside code base, I'll have to introduce an interface to wrap around my gorm.DB to enable mocking for unit test, and change my existing code that directly uses *gorm.DB to instead rely on the interface abstraction. Which is kind of against the my initial intention of using ORM for database manipulation.

However, this is the trade-off of decoupling for testability. Introducing an interface will allow the project no longer tied to a specific ORM (like GORM) in logic layer:

+ Mocking support for unit tests.

+ Cleaner separation of concerns.

+ Easier future refactoring (e.g., switching ORM or DB engine).

For future refactoring (if needed)

- Start by doing this only in controller layer, not everywhere.

- Use GORM directly in low-level repository implementations.

- Gradually migrate more business logic behind interfaces if needed.

### tl;dr:

#### This project uses integration testing instead of unit testing for controller logic.

+ Direct use of *gorm.DB limits testability (cannot mock).

+ Trade-off: simplicity vs decoupling.

#### Suggested Refactor (future):

+ Introduce a repository interface to abstract gorm.DB.

+ Allows mocking and improves separation of concerns.

+ Apply only in controllers first.

#### Benefits of decoupling with interfaces:

+ Unit test support.

+ Cleaner design and future-proofing.

+ Easier to swap ORM/database engines later.

## API document ##
Click [here](./api.md) to API document 