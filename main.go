package main

import (
	"fmt"

	"github.com/BooperDCSS/aggreGator/internal/config"
)

func main() {
	theConfig, err := config.Read()
	if err != nil {
		fmt.Printf("Error reading config file: %v", err)
	}

	theConfig.SetUser("booper")

	updatedConfig, err := config.Read()
	if err != nil {
		fmt.Printf("Error reasding updated config: %v", err)
	}
	fmt.Printf("Current DB URL: %s -- Current Username: %s\n", updatedConfig.Db_url, updatedConfig.Current_Username)

}
