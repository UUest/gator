package commands

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"

	"github.com/UUest/gator/internal/config"
	"github.com/UUest/gator/internal/database"
)

type State struct {
	Config *config.Config
	DB     *database.Queries
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	Names map[string]func(*State, Command) error
}

func HandlerLogin(s *State, cmd Command) error {
	if cmd.Args == nil {
		return fmt.Errorf("Expected a username")
	}
	_, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err == sql.ErrNoRows {
		fmt.Println("User does not exist, please register first")
		os.Exit(1)
	}
	if err != nil {
		return err
	}
	err = s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %s\n", cmd.Args[0])
	return nil
}

func HandlerRegister(s *State, cmd Command) error {
	if cmd.Args == nil {
		return fmt.Errorf("Expected a username")
	}
	existingUser, err := s.DB.GetUser(context.Background(), cmd.Args[0])
	if err != nil && err != sql.ErrNoRows {
		return err
	}
	if err != nil && existingUser.Name == cmd.Args[0] {
		fmt.Println("User already exists")
		os.Exit(1)
	}
	userParams := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
	}
	dbUser, err := s.DB.CreateUser(context.Background(), userParams)
	if err != nil {
		return err
	}
	fmt.Printf("User %s registered\n", dbUser.Name)
	fmt.Printf("ID: %s\n", dbUser.ID)
	fmt.Printf("Created at: %s\n", dbUser.CreatedAt)
	fmt.Printf("Updated at: %s\n", dbUser.UpdatedAt)
	fmt.Printf("Name: %s\n", dbUser.Name)
	HandlerLogin(s, cmd)
	return nil
}

func HandlerReset(s *State, cmd Command) error {
	if err := s.DB.Reset(context.Background()); err != nil {
		return err
	}
	fmt.Println("Database reset")
	return nil
}

func HandlerGetUsers(s *State, cmd Command) error {
	users, err := s.DB.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user.Name == s.Config.CurrentUserName {
			fmt.Printf("* %s (current)\n", user.Name)
		} else {
			fmt.Printf("* %s\n", user.Name)
		}
	}
	return nil
}

func (c *Commands) Run(s *State, cmd Command) error {
	if cmd.Name == "" {
		return fmt.Errorf("No command specified")
	}
	handler, ok := c.Names[cmd.Name]
	if !ok {
		return fmt.Errorf("Unknown command: %s", cmd.Name)
	}
	return handler(s, cmd)
}

func (c *Commands) Register(name string, f func(*State, Command) error) {
	c.Names[name] = f
}
