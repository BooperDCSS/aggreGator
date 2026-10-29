package main

import (
	"fmt"
	"log"
	"os"

	"github.com/BooperDCSS/aggreGator/internal/config"
)

func main() {
	theConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	theState := &state{
		conf: theConfig,
	}

	theCommands := &commands{
		commandList: make(map[string]func(*state, command) error),
	}
	theCommands.register("login", handlerLogin)

	theArgs := os.Args
	if len(theArgs) < 2 {
		log.Fatalf("Too few arguments provided by user; must supply at least a command name.")
	}

	commandInstance := command{
		name:      theArgs[1],
		arguments: theArgs[2:],
	}

	err = theCommands.run(theState, commandInstance)
	if err != nil {
		fmt.Printf("Encountered an error: %v\n", err)
		os.Exit(1)
	}

}
