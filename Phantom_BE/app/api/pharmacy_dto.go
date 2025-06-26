package api

import(
	"PhantomBE/global"
	"time"
)

// 1. Request structure for open pharmacies query
type OpenPharmaciesRequest struct {
    Day  string `json:"day" binding:"required,valid_day" validate_msg:"Day must be a valid day of the week"`
    Time string `json:"time" binding:"required,time_format" validate_msg:"Time must be in HH:MM format"`
}

// 2. Request structure for pharmacies Masks
type PharmacyMasksRequest struct {
	PharmacyID uint   `json:"pharmacy_id" binding:"required,positive_uint" validate_msg:"Pharmacy ID is required and must be greater than 0"`
    Sort       string `json:"sort,omitempty" binding:"omitempty,valid_sort" validate_msg:"Time must be name or price"`
    Order      string `json:"order,omitempty" binding:"omitempty,valid_order" validate_msg:"Time must be asc or desc"`
}
// 3. Request structure for pharmacies filter
type PharmacyFilterRequest struct {
	Operator string  `json:"operator" binding:"required,valid_operator" validate_msg:"Operator is required and must be 'more' or 'less'"`
    Count    int     `json:"count" binding:"required,non_negative_int" validate_msg:"Count is required and cannot be negative"`
    MinPrice float64 `json:"min_price" binding:"required,non_negative_float" validate_msg:"Min price is required and cannot be negative or 0.0"`
    MaxPrice float64 `json:"max_price" binding:"required,non_negative_float" validate_msg:"Max price is required and cannot be negative"`
}
// 4. Request structure for finding top users
type TopUsersRequest struct {
    StartDate string `json:"start_date" binding:"required,date_format" validate_msg:"Start date must be in YYYY-MM-DD format"`
    EndDate   string `json:"end_date" binding:"required,date_format" validate_msg:"End date must be in YYYY-MM-DD format and must be after start date"`
    Limit     int    `json:"limit" binding:"non_negative_int" validate_msg:"Limit must be a positive number or zero"`
}
// 5. Request structure for transaction summary
type TransactionSummaryRequest struct {
	StartDate string `json:"start_date" binding:"required,date_format" validate_msg:"Start date must be in YYYY-MM-DD format"`
    EndDate   string `json:"end_date" binding:"required,date_format" validate_msg:"End date must be in YYYY-MM-DD format and must be after start date"`
}

// 6. Request structure for search
type SearchRequest struct {
	Query string `json:"query" binding:"required,min_search_length,max_search_length,safe_search" validate_msg:"Query must be 2-100 characters and contain only letters, numbers, spaces, hyphens, apostrophes, and periods"`
	Type  string `json:"type" binding:"omitempty,search_type" validate_msg:"Type must be 'pharmacy', 'mask', 'user', or 'all'"`
}

// 7. Request structure for purchase transaction
type PurchaseRequest struct {
    UserID     uint `json:"user_id" binding:"required,positive_uint" validate_msg:"User ID must be a positive number"`
    PharmacyID uint `json:"pharmacy_id" binding:"required,positive_uint" validate_msg:"Pharmacy ID must be a positive number"`
    MaskID     uint `json:"mask_id" binding:"required,positive_uint" validate_msg:"Mask ID must be a positive number"`
    Quantity   int  `json:"quantity" binding:"required,positive_int,max_quantity" validate_msg:"Quantity must be between 1 and 1000"`
}

// Response structure

// 1. Open Pharmacies Response
type OpenPharmaciesResponse struct {
	Pharmacies []global.Pharmacy   `json:"pharmacies"`
	Count      int                 `json:"count"`
}

type OpenPharmacy struct{
	PharmacyID   uint           `json:"pharmacy_id"`
	PharmacyName string         `json:"pharmacy_name"`
}
// 2. Pharmacy Masks Response
type PharmacyMasksResponse struct {
	PharmacyID   uint           `json:"pharmacy_id"`
	PharmacyName string         `json:"pharmacy_name"`
	Masks        []global.Mask  `json:"masks"`
	Count        int            `json:"count"`
}
// 3. Pharmacy Filter Response
type PharmacyWithCount struct {
	global.Pharmacy
	MaskCount int64 `json:"mask_count"`
}
type PharmacyFilterResponse struct{
	Pharmacies []PharmacyWithCount  `json:"pharmacies"`
	Count      int                  `json:"count"`
}

// 4. Top Users Response
type TopUsersResponse struct {
	TopUsers []UserTransactionSummary `json:"top_users"`
	Count    int                      `json:"count"`
	Limit    int                      `json:"limit"`
}

type UserTransactionSummary struct {
	UserID           uint    `json:"user_id"`
	UserName         string  `json:"user_name"`
	TotalAmount      float64 `json:"total_amount"`
	TransactionCount int64   `json:"transaction_count"`
	AverageAmount    float64 `json:"average_amount"`
	Rank             int     `json:"rank"`
}


// 5. Transaction Summary Response
type TransactionSummaryResponse struct {
	Summary TransactionSummaryData `json:"summary"`
}

type TransactionSummaryData struct {
	TotalMasks       int64   `json:"total_masks"`
	TotalValue       float64 `json:"total_value"`
	TransactionCount int64   `json:"transaction_count"`
	AverageValue     float64 `json:"average_value"`
	DailyAverage     float64 `json:"daily_average"`
}

// 6. Search Response
type SearchResponse struct {
	Results []SearchResult `json:"results"`
	Count   int            `json:"count"`
	Query   string         `json:"query"`
	Type    string         `json:"type"`
}

type SearchResult struct {
	Type       string   `json:"type"`
	ID         uint     `json:"id"`
	Name       string   `json:"name"`
	Price      *float64 `json:"price,omitempty"`      // Only for masks
	PharmacyID *uint    `json:"pharmacy_id,omitempty"` // Only for masks
	Relevance  float64  `json:"relevance"`
}

// 7. Purchase Response
type PurchaseResponse struct {
	Success     bool            `json:"success"`
	Message     string          `json:"message"`
	PurchaseIDs []uint          `json:"purchase_ids"`
	Details     PurchaseDetails `json:"details"`
	Timestamp   time.Time       `json:"timestamp"`
}

type PurchaseDetails struct {
	UserID          uint    `json:"user_id"`
	UserName        string  `json:"user_name"`
	PharmacyID      uint    `json:"pharmacy_id"`
	PharmacyName    string  `json:"pharmacy_name"`
	MaskID          uint    `json:"mask_id"`
	MaskName        string  `json:"mask_name"`
	UnitPrice       float64 `json:"unit_price"`
	Quantity        int     `json:"quantity"`
	TotalAmount     float64 `json:"total_amount"`
	PreviousBalance float64 `json:"previous_balance"`
	NewBalance      float64 `json:"new_balance"`
}

// 8. Health check response
type HealthCheckResponse struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
	Reason    string `json:"reason,omitempty"`
	Details   gin.H  `json:"details,omitempty"`
}