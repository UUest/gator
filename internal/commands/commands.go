package commands

import (
	"fmt"

	"github.com/UUest/gator/internal/config"
)

type State struct {
	Config *config.Config
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
	err := s.Config.SetUser(cmd.Args[0])
	if err != nil {
		return err
	}
	fmt.Printf("Username set to %s\n", cmd.Args[0])
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
