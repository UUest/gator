package main

import (
	"fmt"

	"github.com/UUest/gator/internal/config"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	cfg.SetUser("West")
	cfg, err = config.Read()
	if err != nil {
		fmt.Println("Error reading config:", err)
		return
	}
	fmt.Printf("User: %s\n", cfg.CurrentUserName)
	fmt.Printf("DB: %s\n", cfg.DbURL)
}
