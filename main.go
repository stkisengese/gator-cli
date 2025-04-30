package main

import (
	"fmt"
	"log"

	"github.com/stkisengese/gator-cli/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config: %v", err)
	}

	fmt.Println("Initial config:")
	printConfig(cfg)

	yourName := "Neema"
	if err := cfg.SetUser(yourName); err != nil {
		log.Fatalf("Error setting user: %v", err)
	}
	fmt.Printf("\nSuccessfully updated user to: %s\n", yourName)

	// 3. Read config again and print
	updatedCfg, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}
	fmt.Println("\nUpdated config:")
	printConfig(updatedCfg)
}

// Helper function to print config nicely
func printConfig(cfg *config.Config) {
	fmt.Printf("DB URL: %s\n", cfg.DBURL)
	fmt.Printf("Current User: %s\n", cfg.CurrentUserName)
}