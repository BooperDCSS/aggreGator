package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/BooperDCSS/aggreGator/internal/config"
	"github.com/BooperDCSS/aggreGator/internal/database"
	_ "github.com/lib/pq" // using it for side effects, not directly
)

func main() {
	theConfig, err := config.Read()
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	db, err := sql.Open("postgres", theConfig.Db_url)
	dbQueries := database.New(db)
	defer db.Close()

	programState := &state{
		db:   dbQueries,
		conf: theConfig,
	}

	allCommands := &commands{
		commandList: make(map[string]func(*state, command) error),
	}
	allCommands.register("login", handlerLogin)
	allCommands.register("register", handlerRegister)
	allCommands.register("reset", handlerReset)
	allCommands.register("users", handlerListUsers)
	allCommands.register("agg", handlerFetch)
	allCommands.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	allCommands.register("feeds", handlerListFeeds)
	allCommands.register("follow", middlewareLoggedIn(handlerFollowFeed))
	allCommands.register("following", middlewareLoggedIn(handlerCurrentUserFeeds))
	allCommands.register("unfollow", middlewareLoggedIn(handlerDeleteFeedFollow))

	if len(os.Args) < 2 {
		log.Fatal("Too few arguments provided by user; must supply at least a command name.")
	}

	cmdName := os.Args[1]
	cmdArgs := os.Args[2:] // as long as len(os.Args) >= 2, this resolves appropriately

	// using a struct literal for the command argument
	err = allCommands.run(programState, command{name: cmdName, arguments: cmdArgs})
	if err != nil {
		log.Fatalf("Encountered error: %v", err) // this produces exit status code 1
	}

}
