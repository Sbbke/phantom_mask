package global

import (
	"os"
	"time"
)


type OpeningHour struct {
	ID         uint   `gorm:"primaryKey"`
	PharmacyID uint   `gorm:"index"` // foreign key
	DayOfWeek string  `json:"day"`
	OpenTime  string  `json:"open"`  // "HH:MM"
	CloseTime string  `json:"close"` // "HH:MM"
}

type RawPurchase struct {
	PharmacyName      string `json:"pharmacyName"`
	MaskName          string `json:"maskName"`
	TransactionAmount float64 `json:"transactionAmount"`
	TransactionDate   string `json:"transactionDate"` // string from JSON
}

type Purchase struct {
	ID                uint    `gorm:"primaryKey"`
	UserID            uint
	PharmacyName      string  `json:"pharmacyName"`
	MaskName          string  `json:"maskName"`
	TransactionAmount float64 `json:"transactionAmount"`
	TransactionDate   time.Time  `json:"transactionDate"`  
}

type RawUser struct {
	Name              string        `json:"name"`
	CashBalance       float64       `json:"cashBalance"`
	PurchaseHistories []RawPurchase `json:"purchaseHistories"`
}

type User struct {
	ID                uint       `gorm:"primaryKey"`
	Name              string     `json:"name"`
	CashBalance       float64    `json:"cashBalance"`
	PurchaseHistories []Purchase `gorm:"foreignKey:UserID" json:"purchaseHistories"`
}

type Mask struct {
	ID         uint    `gorm:"primaryKey"`
	Name       string  `json:"name"`
	Price      float64 `json:"price"`
	PharmacyID uint
}

type PharmacyRaw struct {
	Name            string         `json:"name"`
	CashBalance     float64        `json:"cashBalance"`
	OpeningHoursRaw string         `json:"openingHours"`
	Masks           []Mask         `json:"masks"`
}

type Pharmacy struct {
	ID              uint           `gorm:"primaryKey"`
	Name            string         `json:"name"`
	CashBalance     float64        `json:"cashBalance"`
	OpeningHours    []OpeningHour  `json:"openingHours" gorm:"foreignKey:PharmacyID"`
	Masks           []Mask         `json:"masks" gorm:"foreignKey:PharmacyID"`
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Error Response
type ErrorResponse struct {
	Error   string      `json:"error"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

var(

	// PostgresUSER defines the configuration for the USER database
	PostgresUser = DatabaseConfig{
		Host:     getEnv("DB_USER_HOST", "postgres-user"),
		Port:     getEnv("DB_USER_PORT", "5432"),
		User:     getEnv("DB_USER_USER", ""),
		Password: getEnv("DB_USER_PASSWORD", ""),
		DBName:   getEnv("DB_USER_DBNAME", "user"),
	}

	// PostgresPHARMACY defines the configuration for the PHARMACY database
	PostgresPharmacy = DatabaseConfig{
		Host:     getEnv("DB_PHARMACY_HOST", "postgres-pharmacy"),
		Port:     getEnv("DB_PHARMACY_PORT", "5432"),
		User:     getEnv("DB_PHARMACY_USER", ""),
		Password: getEnv("DB_PHARMACY_PASSWORD", ""),
		DBName:   getEnv("DB_PHARMACY_DBNAME", "pharmacy"),
	}
	// gin addr
    GinAddr = getEnv("GIN_DOMAIN","") + ":" + getEnv("GIN_PORT", "8080")
	PostgresUserSampleDataFile = getEnv("USER_SAMPLE_FILE", "user.json")
    PostgresPharmacySampleDataFile = getEnv("PHARMACY_SAMPLE_FILE", "pharmacies.json")

	Days = []string{"Monday", "Tuesday", "Wednesday", "Thursday", "Friday", "Saturday", "Sunday"}
	ShortToFullDay = map[string]string{
		"Mon": "Monday", "Tue": "Tuesday", "Wed": "Wednesday",
		"Thu": "Thursday", "Fri": "Friday", "Sat": "Saturday", "Sun": "Sunday",
	}
)

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}