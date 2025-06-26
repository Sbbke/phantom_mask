package validation

import (
	"PhantomBE/app/api"
	"PhantomBE/global"
	"fmt"
	"time"
	"strings"
	"unicode"
	"reflect"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)


// RegisterCustomValidators registers all custom validators used across the app.
func RegisterPharmacyValidators() error {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		// --- Time and Date Validators ---

		// Valid day of the week (e.g., "Monday")
        if err := v.RegisterValidation("valid_day", func(fl validator.FieldLevel) bool {
            _, valid := isValidDay(fl.Field().String())
            return valid
        }); err != nil {
            panic(fmt.Sprintf("Failed to register valid_day validator: %v", err))
        }

		// Time in HH:MM format
        if err := v.RegisterValidation("time_format", func(fl validator.FieldLevel) bool {
            _, err := time.Parse("15:04", fl.Field().String())
            return err == nil
        }); err != nil {
            panic(fmt.Sprintf("Failed to register time_format validator: %v", err))
        }

		// Date in YYYY-MM-DD format
		if err := v.RegisterValidation("date_format", func(fl validator.FieldLevel) bool {
			_, err := time.Parse("2006-01-02", fl.Field().String())
			return err == nil
		}); err != nil {
			panic(fmt.Sprintf("Failed to register date_format validator: %v", err))
		}

		// --- Numeric Validators ---

		// Positive integers
        if err := v.RegisterValidation("positive_int", func(fl validator.FieldLevel) bool {
            return fl.Field().Int() > 0
        }); err != nil {
            panic(fmt.Sprintf("Failed to register non_negative_int validator: %v", err))
        }

		// Non-negative integers (>= 0)
        if err := v.RegisterValidation("non_negative_int", func(fl validator.FieldLevel) bool {
            return fl.Field().Int() >= 0
        }); err != nil {
            panic(fmt.Sprintf("Failed to register non_negative_int validator: %v", err))
        }

		// Positive unsigned integers
        if err := v.RegisterValidation("positive_uint", func(fl validator.FieldLevel) bool {
            return fl.Field().Uint() > 0
        }); err != nil {
            panic(fmt.Sprintf("Failed to register positive_uint validator: %v", err))
        }

		// Positive float
        if err := v.RegisterValidation("non_negative_float", func(fl validator.FieldLevel) bool {
            return fl.Field().Float() > 0
        }); err != nil {
            panic(fmt.Sprintf("Failed to register non_negative_float validator: %v", err))
        }

		// Quantity < 1000
		if err := v.RegisterValidation("max_quantity", func(fl validator.FieldLevel) bool {
			value := fl.Field().Int()
			return value < 1000 // Reasonable maximum to prevent abuse
		}); err != nil {
			panic(fmt.Sprintf("Failed to register max_quantity validator: %v", err))
		}

		// --- Sorting & Filtering ---

		// Valid sort category
		if err := v.RegisterValidation("valid_sort", func(fl validator.FieldLevel) bool {
            if fl.Field().String() == "" {
                return true // Allow empty, will use default
            }
            validSorts := []string{"name", "price"}
            return contains(validSorts, fl.Field().String())
        }); err != nil {
            panic(fmt.Sprintf("Failed to register valid_sort validator: %v", err))
        }

		// Valid order
        if err := v.RegisterValidation("valid_order", func(fl validator.FieldLevel) bool {
            if fl.Field().String() == "" {
                return true // Allow empty, will use default
            }
            validOrders := []string{"asc", "desc"}
            return contains(validOrders, fl.Field().String())
        }); err != nil {
            panic(fmt.Sprintf("Failed to register valid_order validator: %v", err))
        }

		// Valid filter operator
        if err := v.RegisterValidation("valid_operator", func(fl validator.FieldLevel) bool {
            validOperators := []string{"more", "less"}
            return contains(validOperators, fl.Field().String())
        }); err != nil {
            panic(fmt.Sprintf("Failed to register valid_operator validator: %v", err))
        }
		
		// --- Search Query Validators ---

		// Min query length
		if err := v.RegisterValidation("min_search_length", func(fl validator.FieldLevel) bool {
			query := strings.TrimSpace(fl.Field().String())
			return len(query) >= 2 // Minimum 2 characters for meaningful search
		}); err != nil {
			panic(fmt.Sprintf("Failed to register min_search_length validator: %v", err))
		}

		// Max query length
		if err := v.RegisterValidation("max_search_length", func(fl validator.FieldLevel) bool {
			query := strings.TrimSpace(fl.Field().String())
			return len(query) <= 100 // Maximum 100 characters to prevent abuse
		}); err != nil {
			panic(fmt.Sprintf("Failed to register max_search_length validator: %v", err))
		}

		// Valid search type
		if err := v.RegisterValidation("search_type", func(fl validator.FieldLevel) bool {
			searchType := strings.ToLower(strings.TrimSpace(fl.Field().String()))
			validTypes := []string{"pharmacy", "mask", "user", "all"}
			
			for _, validType := range validTypes {
				if searchType == validType {
					return true
				}
			}
			return false
		}); err != nil {
			panic(fmt.Sprintf("Failed to register search_type validator: %v", err))
		}

		// Allow only safe characters in search
		if err := v.RegisterValidation("safe_search", func(fl validator.FieldLevel) bool {
			query := fl.Field().String()
			// Allow alphanumeric, spaces, hyphens, apostrophes, and periods
			for _, char := range query {
				if !unicode.IsLetter(char) && !unicode.IsDigit(char) &&
				char != ' ' && char != '-' && char != '\'' && char != '.' {
					return false
				}
			}
			return true
		}); err != nil {
			panic(fmt.Sprintf("Failed to register safe_search validator: %v", err))
		}

		// --- Struct-Level Validators ---

		// StartDate <= EndDate and within max duration
		dateRangeValidator := func(sl validator.StructLevel) {
			// Use reflection to get StartDate and EndDate fields
			structValue := sl.Current()
			startDateField := structValue.FieldByName("StartDate")
			endDateField := structValue.FieldByName("EndDate")
			
			if !startDateField.IsValid() || !endDateField.IsValid() {
				return // Fields don't exist
			}
			
			startDateStr := startDateField.String()
			endDateStr := endDateField.String()
			
			if startDateStr == "" || endDateStr == "" {
				return // Let individual field validators handle empty values
			}

			startDate, err1 := time.Parse("2006-01-02", startDateStr)
			endDate, err2 := time.Parse("2006-01-02", endDateStr)
			
			if err1 != nil || err2 != nil {
				return // Let individual field validators handle format errors
			}

			// Check if start date is after end date
			if startDate.After(endDate) {
				sl.ReportError(endDateStr, "EndDate", "EndDate", "date_range", "")
			}

			// Check if date range is more than 365 days (configurable)
			maxDuration := 365 * 24 * time.Hour
			if endDate.Sub(startDate) > maxDuration {
				sl.ReportError(endDateStr, "EndDate", "EndDate", "max_date_range", "")
			}
		}

		// Register the combined validation for both structs
		v.RegisterStructValidation(dateRangeValidator, api.TopUsersRequest{})
		v.RegisterStructValidation(dateRangeValidator, api.TransactionSummaryRequest{})
    }
	return nil
}

// formator for reflecting validate_msg
func FormatValidationError(err error, obj interface{}) map[string]string {
	errorsMap := make(map[string]string)

	if validationErrors, ok := err.(validator.ValidationErrors); ok {
		t := reflect.TypeOf(obj)

		for _, fieldErr := range validationErrors {
			fieldName := fieldErr.StructField()

			if f, found := t.FieldByName(fieldName); found {
				msg := f.Tag.Get("validate_msg")
				if msg != "" {
					errorsMap[fieldErr.Field()] = msg
				} else {
					errorsMap[fieldErr.Field()] = fieldErr.Error()
				}
			}
		}
	}

	return errorsMap
}

// Helper function to validate day of week using global Days slice
func isValidDay(day string) (string, bool) {
	for _, validDay := range global.Days {
		if strings.EqualFold(day, validDay) {
			return validDay, true // Return the properly cased day name
		}
	}
	return "", false
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}