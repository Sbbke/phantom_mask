// File: app/middleware/recovery.go
package middleware

import (
	"PhantomBE/global"
	"log"
	"net/http"
	"runtime/debug"
	"time"
	"context"
	"github.com/gin-gonic/gin"
)

// RecoveryMiddleware handles panics and converts them to proper HTTP responses
func RecoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		// Log the panic with stack trace for debugging
		log.Printf("Panic recovered in %s %s: %v\nStack trace:\n%s", 
			c.Request.Method, 
			c.Request.URL.Path, 
			recovered, 
			debug.Stack())
		
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Internal server error",
		})
		c.Abort()
	})
}

// DatabaseCheckMiddleware ensures the database is available
func DatabaseCheckMiddleware(db interface{}) gin.HandlerFunc {
    return func(c *gin.Context) {
        if db == nil {
            c.JSON(http.StatusInternalServerError, global.ErrorResponse{
                Error: "Database connection unavailable",
                Code:  "DB_UNAVAILABLE",
            })
            c.Abort()
            return
        }
        c.Set("db", db)
        c.Next()
    }
}

// DatabaseErrorMiddleware handles database errors from context
func DatabaseErrorMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        c.Next()
        if val := c.Request.Context().Value(global.DBErrorKey); val != nil {
            if err, ok := val.(error); ok {
                c.JSON(http.StatusInternalServerError, global.ErrorResponse{
                    Error: "Database error",
                    Code:  "DB_ERROR",
                    Details: gin.H{
                        "error": err.Error(),
                    },
                })
                c.Abort()
            }
        }
        if val := c.Request.Context().Value(global.DBTimeoutKey); val != nil {
            if err, ok := val.(error); ok {
                c.JSON(http.StatusRequestTimeout, global.ErrorResponse{
                    Error: "Database query timeout",
                    Code:  "DB_TIMEOUT",
                    Details: gin.H{
                        "error": err.Error(),
                    },
                })
                c.Abort()
            }
        }
    }
}

// TimeoutMiddleware enforces a request timeout
func TimeoutMiddleware(timeout time.Duration) gin.HandlerFunc {
    return func(c *gin.Context) {
        ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)
        defer cancel()
        c.Request = c.Request.WithContext(ctx)
        c.Next()
    }
}