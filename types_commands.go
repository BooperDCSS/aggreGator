package main

import (
	"github.com/BooperDCSS/aggreGator/internal/config"
	"github.com/BooperDCSS/aggreGator/internal/database"
)

type state struct {
	db   *database.Queries
	conf *config.Config
}

type command struct {
	name      string
	arguments []string
}

type commands struct {
	commandList map[string]func(*state, command) error
}
