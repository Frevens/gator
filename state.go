package main

import (
    "github.com/Frevens/gator/internal/config"
    "github.com/Frevens/gator/internal/database"
)


type state struct {
	db  *database.Queries
	cfg *config.Config
}