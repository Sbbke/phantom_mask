
package integration

import (
    "PhantomBE/app/models"
    "PhantomBE/app/controllers"
    "PhantomBE/app/api"
    "PhantomBE/app/middleware"
    "PhantomBE/app/validation"
    // "PhantomBE/global"
    "bytes"
    "context"
    "encoding/json"
    // "errors"
    "net/http"
    "net/http/httptest"
    "testing"
    "time"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
    "github.com/charmbracelet/log"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)


func setupRouter(pc *controllers.PharmacyController) *gin.Engine {
    gin.SetMode(gin.TestMode)
    router := gin.Default()
    router.Use(middleware.AddCommonHeaders())
	router.Use(middleware.RecoveryMiddleware())
	router.Use(middleware.DatabaseErrorMiddleware())
	router.Use(middleware.TimeoutMiddleware(60*time.Second))
    router.POST("/api/pharmacies/open", pc.GetOpenPharmacies)
    return router
}

func TestGetOpenPharmacies(t *testing.T) {
    models.ConnectToDatabases("PHARMACY", "USER")
    defer models.CloseConnects( "PHARMACY","USER")
    if err := validation.RegisterPharmacyValidators(); err != nil {
		log.Fatalf("Validator registration failed: %v", err)
	}

    request := api.OpenPharmaciesRequest{
            Day:  "Monday",
            Time: "10:00",
        }
    body, _ := json.Marshal(request)
    t.Run("ValidRequest", func(t *testing.T) {

        pc := controllers.NewPharmacyController(models.DBPharmacy)

        router := setupRouter(pc)
        
        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        var resp api.OpenPharmaciesResponse
        err := json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NoError(t, err)
        assert.Equal(t, 15, resp.Count)
        assert.Len(t, resp.Pharmacies, 15)
        assert.Equal(t, "DFW Wellness", resp.Pharmacies[0].Name)
    })

    t.Run("NilDatabaseConnection", func(t *testing.T) {
        badDB, _ := gorm.Open(postgres.Open("host=invalid port=5432 user=bad dbname=bad"), &gorm.Config{})

        pc := controllers.NewPharmacyController(badDB)
        router := setupRouter(pc)

        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open",  bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()

        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusInternalServerError, w.Code)
        assert.Contains(t, w.Body.String(), "Database error")
    })

    t.Run("InvalidJSON", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBPharmacy)

        router := setupRouter(pc)
		body := bytes.NewBufferString(`{"day":"Monday","time":"10:00"`) // Missing closing brace
		req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", body)
		req.Header.Set("Content-Type", "application/json") // Corrected line
		w := httptest.NewRecorder()
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		var resp gin.H
		err := json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NoError(t, err)
		assert.Contains(t, resp["error"], "Invalid request format")
    })

    t.Run("InvalidDay", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBPharmacy)

        router := setupRouter(pc)
        request := api.OpenPharmaciesRequest{
            Day:  "InvalidDay",
            Time: "10:00",
        }
        body, _ := json.Marshal(request)
        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)

        assert.Contains(t, w.Body.String(), "Day must be a valid day")

    })

    t.Run("InvalidTimeFormat", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBPharmacy)

        router := setupRouter(pc)
        request := api.OpenPharmaciesRequest{
            Day:  "Monday",
            Time: "10:00:00",
        }
        body, _ := json.Marshal(request)
        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusBadRequest, w.Code)

        assert.Contains(t, w.Body.String(), "Time must be in HH:MM format")

    })

    t.Run("DatabaseError", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBUser)

        router := setupRouter(pc)
        request := api.OpenPharmaciesRequest{
            Day:  "Monday",
            Time: "10:00",
        }
        body, _ := json.Marshal(request)
        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusInternalServerError, w.Code)

        assert.Contains(t, w.Body.String(), "Database error")

    })

    t.Run("EmptyResultSet", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBPharmacy)
        router := setupRouter(pc)
        request := api.OpenPharmaciesRequest{
            Day:  "Monday",
            Time: "00:00",
        }
        body, _ := json.Marshal(request)
        req, _ := http.NewRequest(http.MethodPost, "/api/pharmacies/open", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")
        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusOK, w.Code)
        var resp api.OpenPharmaciesResponse
        err := json.Unmarshal(w.Body.Bytes(), &resp)
        assert.NoError(t, err)
        assert.Equal(t, 0, resp.Count)
        assert.Empty(t, resp.Pharmacies)
    })

    t.Run("ContextTimeout", func(t *testing.T) {
        pc := controllers.NewPharmacyController(models.DBPharmacy)
        router := setupRouter(pc)

        router.POST("/api/pharmacies/timeout", func(c *gin.Context) {
            ctx, cancel := context.WithTimeout(c.Request.Context(), 1*time.Nanosecond)
            defer cancel()

            // Force timeout
            time.Sleep(10 * time.Millisecond)

            // Replace context
            c.Request = c.Request.WithContext(ctx)
            pc.GetOpenPharmacies(c)
        })

        req, _ := http.NewRequest("POST", "/api/pharmacies/timeout", bytes.NewBuffer(body))
        req.Header.Set("Content-Type", "application/json")

        w := httptest.NewRecorder()
        router.ServeHTTP(w, req)

        assert.Equal(t, http.StatusRequestTimeout, w.Code)
    })
    t.Run("PanicRecovery", func(t *testing.T){
        pc := controllers.NewPharmacyController(models.DBPharmacy)
        router := setupRouter(pc)
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
