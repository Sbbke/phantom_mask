package app


import (

	"PhantomBE/app/routes"
	"PhantomBE/global"
	"PhantomBE/app/models"
	"PhantomBE/app/initial"
	"PhantomBE/app/middleware"
	"PhantomBE/app/validation"
	"time"
	"github.com/charmbracelet/log"
	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

// app.go is the main entry point for this application.
// initiates configures the postgreSQL, validatior, middleware, Gin router and starts the server.

// StartApplication initiates the main backend application, including the Gin router and postgreSQL.
func StartApplication() {
	// 1. Connect to postgreSQL and Redis
	models.ConnectToDatabases("PHARMACY")
	
	// 2. Registor custom validators
	if err := validation.RegisterPharmacyValidators(); err != nil {
		log.Fatalf("Validator registration failed: %v", err)
	}
	// 3. Initiate default Gin router with logger and recovery middleware
	routes.Router = gin.Default()

	// 4. Add common middlewares that need to run on all routes
	routes.Router.Use(middleware.AddCommonHeaders())
	routes.Router.Use(middleware.RecoveryMiddleware())
	routes.Router.Use(middleware.DatabaseErrorMiddleware())
	routes.Router.Use(middleware.TimeoutMiddleware(60*time.Second))
	// 5. Configure routes and routing groups (./router.go)
	routes.ConfigureRoutes()

	// 6. Configure http server
	addr := global.GinAddr

	err := endless.ListenAndServe(addr, routes.Router)
	if err != nil {
		log.Warn(err)
	}
	log.Info("Server on %v stopped", addr)

	// If the server stops, close the database connections
	models.CloseConnects( "PHARMACY")

}

// preprocess data from user.json
func InitUserSchema() {
	models.ConnectToDatabases("PHARMACY")
	initial.InitSampleUser()
	models.CloseConnects("PHARMACY")
}

// preprocess data from pharmacies.json
func InitPharmaciesData() {
	models.ConnectToDatabases("PHARMACY")
	initial.InitSamplePharmacies()
	models.CloseConnects("PHARMACY")
}

// migrate preprocessed data
func MigrateData() {
	models.ConnectToDatabases("PHARMACY")
	if err := models.MigrateSchema(); err != nil {
		log.Fatal("migration failed: %v", err)
		return
	}
	models.CloseConnects("PHARMACY")
}