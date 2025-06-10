package main

import (
	"fmt"
	"os"

	"github.com/UUest/gator/internal/commands"
	"github.com/UUest/gator/internal/config"
	_ "github.com/lib/pq"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		os.Exit(1)
	}
	s := commands.State{
		Config: cfg,
	}

	c := commands.Commands{
		Names: make(map[string]func(*commands.State, commands.Command) error),
	}
	c.Register("login", commands.HandlerLogin)

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
