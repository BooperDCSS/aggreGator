package main

import (
	"log"
	"os"

	"github.com/BooperDCSS/aggreGator/internal/config"
)

func main() {
	theConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	programState := &state{
		conf: theConfig,
	}

	allCommands := &commands{
		commandList: make(map[string]func(*state, command) error),
	}
	allCommands.register("login", handlerLogin)

	if len(os.Args) < 2 {
		log.Fatal("Too few arguments provided by user; must supply at least a command name.")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:]

	// using a struct literal for the command argument
	err = allCommands.run(programState, command{name: cmdName, arguments: cmdArgs})
	if err != nil {
		log.Fatal(err) // this produces exit status code 1
	}

}
