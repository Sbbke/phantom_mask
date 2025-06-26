
package integration

import (

    "PhantomBE/app/middleware"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"

)


func setupRouter() *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.Use(middleware.AddCommonHeaders())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.DatabaseErrorMiddleware())
	router.Use(middleware.TimeoutMiddleware(60*time.Second))
    return router
}

func TestMiddleware(t *testing.T) {

    t.Run("PanicRecovery", func(t *testing.T){

        router := setupRouter()
        router.POST("/test/panic", func(c *gin.Context) {
            panic("simulated panic for testing")
        })

        req, _ := http.NewRequest("POST", "/test/panic", nil)
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusInternalServerError, w.Code)
        assert.Contains(t, w.Body.String(), "Internal server error")
    })
   
}
