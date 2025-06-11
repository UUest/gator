package main

import (
	"fmt"
	"os"

	"database/sql"

	"github.com/UUest/gator/internal/commands"
	"github.com/UUest/gator/internal/config"
	"github.com/UUest/gator/internal/database"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		fmt.Println("Error opening database:", err)
		os.Exit(1)
	}
	defer db.Close()
	dbQueries := database.New(db)

	s := commands.State{
		Config: cfg,
		DB:     dbQueries,
	}

	c := commands.Commands{
		Names: make(map[string]func(*commands.State, commands.Command) error),
	}
	c.Register("login", commands.HandlerLogin)
	c.Register("register", commands.HandlerRegister)

	input := os.Args
	if len(input) < 3 {
		fmt.Println("Usage: gator <command> <args>")
		os.Exit(1)
	}

	commandName := input[1]
	commandArgs := input[2:]

	cmd := commands.Command{
		Name: commandName,
		Args: commandArgs,
	}

	err = c.Run(&s, cmd)
	if err != nil {
		fmt.Println("Error running command:", err)
		os.Exit(1)
	}

}
