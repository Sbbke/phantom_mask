package cmd

import (
	"fmt"
	"os"
	"github.com/charmbracelet/log"
	"PhantomBE/app"
	"github.com/spf13/cobra"
)

// rootCmd represents the root command for the PhantomBE application.
var rootCmd = &cobra.Command{
	Use:   "PhantomBE",
	Short: "Phantom application backend",
	Long:  "Backend application for APIs",
	Run: func(_ *cobra.Command, _ []string) {
		log.Info("Welcome to Phantom Backend!")
		// Start the application when the root command is executed.
		app.StartApplication()
	},
}

// preprocess user data
var initUsersDBCmd = &cobra.Command{
	Use:   "initUsers",
	Short: "init user data",
	Long:  "Preprocess user data",
	Run: func(_ *cobra.Command, _ []string) {
		log.Info("Start the process of init user data.")
		app.InitUserSchema()
	},
}

// preprocess pharmacies data
var initPharmaciesDBCmd = &cobra.Command{
	Use:   "initPharmacies",
	Short: "init Pharmacies data",
	Long:  "Preprocess pharmacies data",
	Run: func(_ *cobra.Command, _ []string) {
		log.Info("Start the process of init pharmacies data.")
		app.InitPharmaciesData()
	},
}

// migrate preprocessed data
var migrateSchemaCMD = &cobra.Command{
	Use:   "migrateSchema",
	Short: "migrate database",
	Long:  "Use models package to Create or Update pharmacy DB table Schema.",
	Run: func(_ *cobra.Command, _ []string) {
		log.Info("Start the process of insert pharmacies database data.")
		app.MigrateData()
	},
}

// Execute initializes Cobra and adds the checkExpiredCmd to the root command.
func Execute() {
	rootCmd.AddCommand(initUsersDBCmd)
	rootCmd.AddCommand(initPharmaciesDBCmd)
	rootCmd.AddCommand(migrateSchemaCMD)
	// Execute the root command and handle any errors.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}