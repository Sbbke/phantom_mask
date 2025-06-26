package models

import (
	// "bufio"
	// "database/sql"
	"fmt"
	// "os"
	// "strings"

	"PhantomBE/global"
	"github.com/charmbracelet/log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DBPharmacy stores the actual statistical data for each component in the Pharmacy.
// DBUser stores the user data and various configs for this application.

// Global Variables that allow access to the DB anywhere in the application.
var (
	DBPharmacy *gorm.DB
	DBUser   *gorm.DB
)

// ConnectToDatabases connects to the two PostgreSQL databases used by this application.
// It accepts a variable number of database names as arguments and establishes connections to each specified database.
// The function iterates over the provided database names and connects to the corresponding databases.
// It uses a switch statement to handle different database names and assigns the database connections accordingly.
func ConnectToDatabases(dbNames ...interface{}) {
	for _, dbName := range dbNames {
		if dbString, ok := dbName.(string); ok {
			// Switch statement to handle different database names.
			switch dbString {
			case "PHARMACY":
				log.Info("PHARMACY ", "Hostname", global.PostgresPharmacy.Host)
				DBPharmacy = ConnectToDatabase(global.PostgresPharmacy)
			case "USER":
				log.Info("USER ","Hostname", global.PostgresUser.Host)
				DBUser = ConnectToDatabase(global.PostgresUser)
			default:
				panic("Database not in connection list.")
			}
		}
	}
}

// ConnectToDatabase establishes a connection to the specified database using the provided DatabaseConfig.
// It constructs the database connection string using the provided database configuration and establishes a connection using gorm.Open.
// The function returns a pointer to a gorm.DB (database connection) and logs success or failure messages accordingly.
func ConnectToDatabase(dbConfig global.DatabaseConfig) *gorm.DB {
	// Constructing the database connection string using database configuration
	dbargs := fmt.Sprintf(
		"host=%s port=%s user=%s dbname=%s password=%s sslmode=disable",
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.User,
		dbConfig.DBName,
		dbConfig.Password,
	)

	// Establish a connection to the database using gorm.Open and the constructed connection string
	dbConn, err := gorm.Open(postgres.Open(dbargs), &gorm.Config{})
	if err != nil {
		// Log an error and panic if there is an issue connecting to the database
		log.Error("Error connecting to database", "database", dbConfig.Host)
		panic("Connecting to database error")
	}

	// Log a success message if the connection is established successfully
	log.Info("database connected", "database",dbConfig.Host)
	return dbConn
}

// CloseConnects closes the connections to the specified databases.
// It takes a variable number of database names and closes the corresponding connections.
func CloseConnects(dbNames ...interface{}) {
	for _, dbName := range dbNames {
		if dbString, ok := dbName.(string); ok {
			// Switch statement to handle different database names.
			switch dbString {
			case "PHARMACY":
				CloseConnect(dbString, DBPharmacy)
			case "USER":
				CloseConnect(dbString, DBUser)
			default:
				panic("DB does not in connection list.")
			}
		}
	}
}

// CloseConnect closes the connection to the specified database.
// It takes the database name and the corresponding *gorm.DB object as parameters.
func CloseConnect(dbName string, DB *gorm.DB) {
	// Retrieve the underlying SQL database connection.
	sqlDB, err := DB.DB()
	if err != nil {
		log.Error("failed to get database connection", "database", dbName)
	}

	// Close the underlying SQL database connection.
	if err := sqlDB.Close(); err != nil {
		log.Error("failed to close database connection", "database", dbName)
	} else {
		log.Info("database connection closed", "database", dbName)
	}
}
func MigrateSchema() error {
	// Retrieve the underlying SQL database connection.
	if err := DBPharmacy.AutoMigrate(&global.User{}, &global.Purchase{}, &global.Pharmacy{}, &global.Mask{}, &global.OpeningHour{}); err != nil {
		log.Error("failed to auto migrate DB", "err" , err)
		return err
	}
	log.Info("Schema migrated successfully")
	return nil
}
