package main

import (
	"fmt"
	"log"

	"github.com/BooperDCSS/aggreGator/internal/config"
)

func main() {
	theConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	theConfig.SetUser("booper")

	theConfig, err = config.Read()
	if err != nil {
		log.Fatalf("Error reading updated config: %v", err)
	}
	fmt.Printf("Updated config: %+v", theConfig)

}
