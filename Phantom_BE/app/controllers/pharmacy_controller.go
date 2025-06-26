package controllers
import (
	"PhantomBE/global"
	"PhantomBE/app/api"
	"PhantomBE/app/validation"
	"gorm.io/gorm"
	"strings"
	"sort"
	"net/http"
	"context"
	"errors"
	// "fmt"
	"time"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type PharmacyController struct {
	db *gorm.DB
}

func NewPharmacyController(db *gorm.DB) *PharmacyController {
	return &PharmacyController{db: db}
}


// 1. List all pharmacies open at a specific time and on a day of the week
// POST /api/v1/pharmacies/open
func (pc *PharmacyController) GetOpenPharmacies(c *gin.Context) {

	// Add timeout context to prevent hanging requests
	ctx := c.Request.Context()


	var req api.OpenPharmaciesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	var pharmacies []global.Pharmacy
	
	err := pc.db.WithContext(ctx).
        Select("pharmacies.*").
        Joins("JOIN opening_hours ON pharmacies.id = opening_hours.pharmacy_id").
        Where("opening_hours.day_of_week = ? AND ? BETWEEN opening_hours.open_time AND opening_hours.close_time", req.Day, req.Time).
        Limit(1000).
        Find(&pharmacies).Error

	if err != nil {
        // Set context for middleware to handle
        key := global.DBErrorKey
        if errors.Is(err, context.DeadlineExceeded) {
            key = global.DBTimeoutKey
        }
        ctx = context.WithValue(ctx, key, err)
        c.Request = c.Request.WithContext(ctx)
		c.Abort()
        return 
    }

	response := api.OpenPharmaciesResponse{
		Pharmacies: pharmacies,
		Count:      len(pharmacies),
	}
	
	c.JSON(http.StatusOK, response)
}
// 2. List all masks sold by a given pharmacy, sorted by mask name or price
// POST /api/v1/pharmacies/masks
func (pc *PharmacyController) GetPharmacyMasks(c *gin.Context) {

		// Add timeout context to prevent hanging requests
	ctx := c.Request.Context()


	var req api.PharmacyMasksRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	// Set defaults
	if req.Sort == "" {
		req.Sort = "name"
	}
	if req.Order == "" {
		req.Order = "asc"
	}

	// Check if pharmacy exists
	var pharmacy global.Pharmacy
	err := pc.db.WithContext(ctx).First(&pharmacy, req.PharmacyID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, global.ErrorResponse{
				Error: "Pharmacy not found",
				Code:  "PHARMACY_NOT_FOUND",
				Details: gin.H{
					"pharmacy_id": req.PharmacyID,
				},
			})
			return
		}
		// Record error in context for middleware
		key := global.DBErrorKey
        if errors.Is(err, context.DeadlineExceeded) {
            key = global.DBTimeoutKey
        }
        ctx = context.WithValue(ctx, key, err)
        c.Request = c.Request.WithContext(ctx)
        c.Abort()
		return
	}
	// Query masks
	var masks []global.Mask
	orderClause := req.Sort + " " + req.Order
	
	err = pc.db.WithContext(ctx).
        Where("pharmacy_id = ?", req.PharmacyID).
        Order(orderClause).
        Limit(1000).
        Find(&masks).Error
	if err != nil {
        // Handle database errors via middleware
        key := global.DBErrorKey
        if errors.Is(err, context.DeadlineExceeded) {
            key = global.DBTimeoutKey
        }
        ctx = context.WithValue(ctx, key, err)
        c.Request = c.Request.WithContext(ctx)
        c.Abort()
        return
    }

	response := api.PharmacyMasksResponse{
		PharmacyID  :pharmacy.ID,
		PharmacyName: pharmacy.Name,
		Masks:         masks,
		Count:         len(masks),
	}

	c.JSON(http.StatusOK, response)
}

// 3. List all pharmacies with more or less than x mask products within a price range
// POST /api/v1/pharmacies/filter
func (pc *PharmacyController) GetPharmaciesByMaskCount(c *gin.Context) {
	// Add timeout context to prevent hanging requests
	ctx := c.Request.Context()

	var req api.PharmacyFilterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
    }

	if req.MinPrice > req.MaxPrice {
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "min_price cannot be greater than max_price",
			Code:  "INVALID_PRICE_RANGE",
			Details: gin.H{
				"min_price": req.MinPrice,
				"max_price": req.MaxPrice,
			},
		})
		return
	}

	var results []api.PharmacyWithCount
	var havingClause string
	
	if req.Operator == "more" {
		havingClause = "COUNT(masks.id) > ?"
	} else {
		havingClause = "COUNT(masks.id) < ?"
	}

	err := pc.db.WithContext(ctx).
        Table("pharmacies").
        Select("pharmacies.*, COUNT(masks.id) as mask_count").
        Joins("LEFT JOIN masks ON pharmacies.id = masks.pharmacy_id AND masks.price BETWEEN ? AND ?", req.MinPrice, req.MaxPrice).
        Group("pharmacies.id").
        Having(havingClause, req.Count).
        Limit(1000).
        Find(&results).Error

	if err != nil {
		// Record error in context for middleware
		key := global.DBErrorKey
        if errors.Is(err, context.DeadlineExceeded) {
            key = global.DBTimeoutKey
        }
        ctx = context.WithValue(ctx, key, err)
        c.Request = c.Request.WithContext(ctx)
        c.Abort()
        return
	}

	response := api.PharmacyFilterResponse{
		Pharmacies:results,
		Count:len(results),
	}
	c.JSON(http.StatusOK, response)
}

// 4. The top x users by total transaction amount of masks within a date range
// POST /api/v1/pharmacies/users/top
func (pc *PharmacyController) GetTopUsers(c *gin.Context) {
	// Add timeout context to prevent hanging requests
	ctx := c.Request.Context()
	var req api.TopUsersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}


	// Set default limit if not provided or zero
	if req.Limit <= 0 {
		req.Limit = 10
	}

	// Parse dates (validation already ensures correct format)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	// End date should include the entire day
	endDate = endDate.Add(24*time.Hour - time.Second)

	type UserWithTotal struct {
		UserID           uint    `json:"user_id"`
		UserName         string  `json:"user_name"`
		TotalAmount      float64 `json:"total_amount"`
		TransactionCount int64   `json:"transaction_count"`
	}

	var topUsers []UserWithTotal
	err := pc.db.WithContext(ctx).
		Table("purchases AS p").
		Joins("JOIN users u ON u.id = p.user_id").
		Select(
			"u.id AS user_id",
			"u.name AS user_name", 
			"SUM(p.transaction_amount) AS total_amount",
			"COUNT(*) AS transaction_count").
		Where("p.transaction_date BETWEEN ? AND ?", startDate, endDate).
		Group("u.id, u.name").
		Order("total_amount DESC").
		Limit(req.Limit).
		Scan(&topUsers).Error

	if err != nil {
		key := global.DBErrorKey
		if errors.Is(err, context.DeadlineExceeded) {
			key = global.DBTimeoutKey
		}
		ctx = context.WithValue(ctx, key, err)
		c.Request = c.Request.WithContext(ctx)
		c.Abort()
		return
	}
	// Transform results to response format
	responseUsers := make([]api.UserTransactionSummary, len(topUsers))
	for i, user := range topUsers {
		avgAmount := 0.0
		if user.TransactionCount > 0 {
			avgAmount = user.TotalAmount / float64(user.TransactionCount)
		}
		
		responseUsers[i] = api.UserTransactionSummary{
			UserID:           user.UserID,
			UserName:         user.UserName,
			TotalAmount:      user.TotalAmount,
			TransactionCount: user.TransactionCount,
			AverageAmount:    avgAmount,
			Rank:             i + 1,
		}
	}
	
	response := api.TopUsersResponse{
		TopUsers: responseUsers,
		Count:    len(responseUsers),
		Limit: req.Limit,
	}
	c.JSON(http.StatusOK, response)
}

// 5. The total number of masks and dollar value of transactions within a date range
// POST /api/v1/pharmacies/transactions/summary
func (pc *PharmacyController) GetTransactionSummary(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.TransactionSummaryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	// Parse dates (validation already ensures correct format)
	startDate, _ := time.Parse("2006-01-02", req.StartDate)
	endDate, _ := time.Parse("2006-01-02", req.EndDate)

	// End date should include the entire day
	endDate = endDate.Add(24*time.Hour - time.Second)

	type TransactionSummary struct {
		TotalMasks       int64   `json:"total_masks"`
		TotalValue       float64 `json:"total_value"`
		TransactionCount int64   `json:"transaction_count"`
		AverageValue     float64 `json:"average_value"`
	}

	var summary TransactionSummary
	err := pc.db.WithContext(ctx).
		Table("purchases").
		Select("COUNT(*) as total_masks, SUM(transaction_amount) as total_value, COUNT(*) as transaction_count, AVG(transaction_amount) as average_value").
		Where("transaction_date BETWEEN ? AND ?", startDate, endDate).
		Scan(&summary).Error

	if err != nil {
		key := global.DBErrorKey
		if errors.Is(err, context.DeadlineExceeded) {
			key = global.DBTimeoutKey
		}
		ctx = context.WithValue(ctx, key, err)
		c.Request = c.Request.WithContext(ctx)
		c.Abort()
		return
	}

	// Calculate daily average
	daysDiff := endDate.Sub(startDate).Hours() / 24
	if daysDiff <= 0 {
		daysDiff = 1
	}

	dailyAverage := 0.0
	if summary.TotalValue > 0 {
		dailyAverage = summary.TotalValue / daysDiff
	}

	response := api.TransactionSummaryResponse{
		Summary: api.TransactionSummaryData{
			TotalMasks:       summary.TotalMasks,
			TotalValue:       summary.TotalValue,
			TransactionCount: summary.TransactionCount,
			AverageValue:     summary.AverageValue,
			DailyAverage:     dailyAverage,
		},
	}
	c.JSON(http.StatusOK, response)
}

// 6. Search for pharmacies or masks by name, ranked by relevance to the search term
// POST /api/v1/pharmacies/search
func (pc *PharmacyController) Search(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.SearchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	// Set default search type and normalize
	if req.Type == "" {
		req.Type = "all"
	}
	req.Type = strings.ToLower(strings.TrimSpace(req.Type))
	req.Query = strings.TrimSpace(req.Query)

	var results []api.SearchResult

	// Search pharmacies
	if req.Type == "pharmacy" || req.Type == "all" {
		pharmacyResults, err := pc.searchPharmacies(ctx, req.Query)
		if err != nil {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
			return
		}
		results = append(results, pharmacyResults...)
	}

	// Search users
	if req.Type == "user" || req.Type == "all" {
		userResults, err := pc.searchUsers(ctx, req.Query)
		if err != nil {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
			return
		}
		results = append(results, userResults...)
	}

	// Search masks
	if req.Type == "mask" || req.Type == "all" {
		maskResults, err := pc.searchMasks(ctx, req.Query)
		if err != nil {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
			return
		}
		results = append(results, maskResults...)
	}

	// Sort by relevance (higher is better) - using sort.Slice for better performance
	sort.Slice(results, func(i, j int) bool {
		return results[i].Relevance > results[j].Relevance
	})

	response := api.SearchResponse{
		Results: results,
		Count:   len(results),
		Query:   req.Query,
		Type:    req.Type,
	}

	c.JSON(http.StatusOK, response)
}

// 7. Process a user purchases a mask from a pharmacy
// POST /api/v1/pharmacies/purchase
func (pc *PharmacyController) ProcessPurchase(c *gin.Context) {
	ctx := c.Request.Context()

	var req api.PurchaseRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		if ve, ok := err.(validator.ValidationErrors); ok {
			c.JSON(http.StatusBadRequest, global.ErrorResponse{
				Error: "Invalid input",
				Code:  "INVALID_INPUT",
				Details: validation.FormatValidationError(ve, req),
			})
			return
		}
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Invalid request format",
			Code:  "INVALID_REQUEST",
			Details: err.Error(),
		})
		return
	}

	// Start transaction
	tx := pc.db.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Get user with row lock to prevent concurrent modifications
	var user global.User
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&user, req.UserID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, global.ErrorResponse{
				Error: "User not found",
				Code:  "USER_NOT_FOUND",
			})
		} else {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
		}
		return
	}

	// Get mask
	var mask global.Mask
	if err := tx.First(&mask, req.MaskID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, global.ErrorResponse{
				Error: "Mask not found",
				Code:  "MASK_NOT_FOUND",
			})
		} else {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
		}
		return
	}

	// Verify mask belongs to the specified pharmacy
	if mask.PharmacyID != req.PharmacyID {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Mask does not belong to specified pharmacy",
			Code:  "MASK_PHARMACY_MISMATCH",
			Details: gin.H{
				"mask_pharmacy_id": mask.PharmacyID,
				"requested_pharmacy_id": req.PharmacyID,
			},
		})
		return
	}

	// Get pharmacy with row lock
	var pharmacy global.Pharmacy
	if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&pharmacy, req.PharmacyID).Error; err != nil {
		tx.Rollback()
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, global.ErrorResponse{
				Error: "Pharmacy not found",
				Code:  "PHARMACY_NOT_FOUND",
			})
		} else {
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
		}
		return
	}

	// Calculate total amount
	totalAmount := mask.Price * float64(req.Quantity)
	previousBalance := user.CashBalance

	// Check if user has sufficient balance
	if user.CashBalance < totalAmount {
		tx.Rollback()
		c.JSON(http.StatusBadRequest, global.ErrorResponse{
			Error: "Insufficient balance",
			Code:  "INSUFFICIENT_BALANCE",
			Details: gin.H{
				"required_amount": totalAmount,
				"current_balance": user.CashBalance,
				"shortage":        totalAmount - user.CashBalance,
			},
		})
		return
	}

	// Update user balance
	user.CashBalance -= totalAmount
	if err := tx.Save(&user).Error; err != nil {
		tx.Rollback()
		key := global.DBErrorKey
		if errors.Is(err, context.DeadlineExceeded) {
			key = global.DBTimeoutKey
		}
		ctx = context.WithValue(ctx, key, err)
		c.Request = c.Request.WithContext(ctx)
		c.Abort()
		return
	}

	// Update pharmacy balance
	pharmacy.CashBalance += totalAmount
	if err := tx.Save(&pharmacy).Error; err != nil {
		tx.Rollback()
		key := global.DBErrorKey
		if errors.Is(err, context.DeadlineExceeded) {
			key = global.DBTimeoutKey
		}
		ctx = context.WithValue(ctx, key, err)
		c.Request = c.Request.WithContext(ctx)
		c.Abort()
		return
	}

	// Create purchase records (one for each quantity)
	purchaseIDs := make([]uint, 0, req.Quantity)
	for i := 0; i < req.Quantity; i++ {
		purchase := global.Purchase{
			UserID:            req.UserID,
			PharmacyName:      pharmacy.Name,
			MaskName:          mask.Name,
			TransactionAmount: mask.Price,
			TransactionDate:   time.Now(),
		}

		if err := tx.Create(&purchase).Error; err != nil {
			tx.Rollback()
			key := global.DBErrorKey
			if errors.Is(err, context.DeadlineExceeded) {
				key = global.DBTimeoutKey
			}
			ctx = context.WithValue(ctx, key, err)
			c.Request = c.Request.WithContext(ctx)
			c.Abort()
			return
		}
		purchaseIDs = append(purchaseIDs, purchase.ID)
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		key := global.DBErrorKey
		if errors.Is(err, context.DeadlineExceeded) {
			key = global.DBTimeoutKey
		}
		ctx = context.WithValue(ctx, key, err)
		c.Request = c.Request.WithContext(ctx)
		c.Abort()
		return
	}
	
	response := api.PurchaseResponse{
		Success:     true,
		Message:     "Purchase completed successfully",
		PurchaseIDs: purchaseIDs,
		Details: api.PurchaseDetails{
			UserID:          req.UserID,
			UserName:        user.Name,
			PharmacyID:      req.PharmacyID,
			PharmacyName:    pharmacy.Name,
			MaskID:          req.MaskID,
			MaskName:        mask.Name,
			UnitPrice:       mask.Price,
			Quantity:        req.Quantity,
			TotalAmount:     totalAmount,
			PreviousBalance: previousBalance,
			NewBalance:      user.CashBalance,
		},
		Timestamp: time.Now(),
	}
	c.JSON(http.StatusOK, response)
}

// Add a health check endpoint to monitor database connectivity
// POST /api/v1/pharmacies/health
func (pc *PharmacyController) HealthCheck(c *gin.Context) {
	if pc.db == nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "reason": "database connection unavailable"})
		return
	}

	// Test database connectivity
	sqlDB, err := pc.db.DB()
	if err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "reason": "cannot access database"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		c.JSON(http.StatusServiceUnavailable, gin.H{"status": "unhealthy", "reason": "database ping failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "healthy"})
}