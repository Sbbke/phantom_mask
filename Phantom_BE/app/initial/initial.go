package initial

import(

	"PhantomBE/global"
	"PhantomBE/app/models"
	"github.com/charmbracelet/log"
	"os"
	"fmt"
	"encoding/json"
	"time"
	// "gorm.io/gorm"
)

func InitSampleUser() {

	// get import file path
	filePath := global.SampleDataDir + global.PostgresUserSampleDataFile
	// outputPath := "/data/users_preprocessed.json"
	
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("failed to read user JSON file","err", err)
	}
	var rawUsers []global.RawUser
	err = json.Unmarshal(data, &rawUsers)
	if err != nil {
		panic(fmt.Sprintf("failed to unmarshal users json: %v", err))
	}

	// 
	// var processedUsers []global.User
	for _, raw  := range rawUsers {
		var purchases []global.Purchase
		for _, rp := range raw.PurchaseHistories {
			parsedTime, err := time.Parse("2006-01-02 15:04:05", rp.TransactionDate)
			if err != nil {
				log.Info(fmt.Sprintf("failed to parse transaction date %s: %v", rp.TransactionDate, err))
				continue
			}
			purchases = append(purchases, global.Purchase{
				PharmacyName:      rp.PharmacyName,
				MaskName:          rp.MaskName,
				TransactionAmount: rp.TransactionAmount,
				TransactionDate:   parsedTime,
			})
		}
		user := global.User{
			Name:              raw.Name,
			CashBalance:       raw.CashBalance,
			PurchaseHistories: purchases,
		}
		// processedUsers = append(processedUsers, user)

		if err := models.DBPharmacy.Where("name = ?", user.Name).FirstOrCreate(&user).Error; err != nil {
			log.Printf("Failed to insert user %s: %v", raw.Name, err)
		} else {
			log.Printf("Inserted user %s", raw.Name)
		}
	}

	// Save preprocessed data to output file
	// outputData, err := json.MarshalIndent(processedUsers, "", "  ")
	// if err != nil {
	// 	log.Error("failed to marshal processed users: %v", err)
	// }
	// if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
	// 	log.Error("failed to write preprocessed users to %s: %v", outputPath, err)
	// }
	// log.Info("Preprocessed users saved", "Path", outputPath)

}


func InitSamplePharmacies() {
	// get import file path
	filePath := global.SampleDataDir + global.PostgresPharmacySampleDataFile
	// outputPath := "/data/pharmacies_preprocessed.json"

	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Error("failed to read pharmacy JSON file","err", err)
	}
	var rawPharmacies []global.PharmacyRaw
	if err := json.Unmarshal(data, &rawPharmacies); err != nil {
		log.Error("failed to unmarshal pharmacy JSON", "err",err)
	}	

	// var processedPharmacies []global.Pharmacy

	for _, rp := range rawPharmacies {
		// Normalize openingHours
		parsed := parseOpeningHours(rp.OpeningHoursRaw)

		pharmacy := global.Pharmacy{
			Name: rp.Name,
			CashBalance: rp.CashBalance,
			OpeningHours: parsed,
			Masks: rp.Masks,
		}

		// processedPharmacies = append(processedPharmacies, pharmacy)
		if err := models.DBPharmacy.Where("name = ?", pharmacy.Name).FirstOrCreate(&pharmacy).Error; err != nil {
			log.Error("failed to insert pharmacy", "pharmacy", pharmacy.Name, "error", err)
		} else {
			log.Info("inserted pharmacy", "pharmacy", pharmacy.Name)
		}
	}

	// Save preprocessed data to output file
	// outputData, err := json.MarshalIndent(processedPharmacies, "", "  ")
	// if err != nil {
	// 	log.Error("failed to marshal processed pharmacies: %v", err)

	// }
	// if err := os.WriteFile(outputPath, outputData, 0644); err != nil {
	// 	log.Error("failed to write preprocessed pharmacies to %s: %v", outputPath, err)

	// }
	// log.Info("Preprocessed pharmacies saved ", "Path", outputPath)

}
