package routes

import (
	"PhantomBE/app/controllers"
	"PhantomBE/app/models"
	"PhantomBE/global"
	"github.com/gin-gonic/gin"
)

// router.go configures all API routes

var (
	Router      *gin.Engine
	RouterGroup *gin.RouterGroup
)

func ConfigureRoutes() {

	Router.Use(gin.Recovery())
	// API routers
	
	RouterGroup = Router.Group("/api/" + global.VERSION)

	// Adding Authentication for security
	// configureHelloRoutes()
	configureHelloRoute()
	// configure test function routing
	configurePharmacyRoutes()
}

func configureHelloRoute(){
	helloRoutes := RouterGroup.Group("/hello")
	{
		helloRoutes.GET("/",controllers.HelloHandler)
	}
}
func configurePharmacyRoutes() {
	// Use the global database variable
	pc := controllers.NewPharmacyController(models.DBPharmacy)
	
	// Create pharmacy routes group
	pharmacyGroup := RouterGroup.Group("/pharmacies")
	{
		// List pharmacies open at specific time/day
		pharmacyGroup.POST("/open", pc.GetOpenPharmacies)
		pharmacyGroup.POST("/masks", pc.GetPharmacyMasks)
		pharmacyGroup.POST("/filter", pc.GetPharmaciesByMaskCount)
		pharmacyGroup.POST("/users/top", pc.GetTopUsers)
		pharmacyGroup.POST("/transactions/summary", pc.GetTransactionSummary)
		pharmacyGroup.POST("/search", pc.Search)
		pharmacyGroup.POST("/purchase", pc.ProcessPurchase)
		pharmacyGroup.GET("/health", pc.HealthCheck)
		
	}
}