package main

import (
	"github.com/BooperDCSS/aggreGator/internal/config"
	"github.com/BooperDCSS/aggreGator/internal/database"
)

type state struct {
	db     *database.Queries
	conf   *config.Config
	target string
}
