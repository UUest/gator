package commands

import (
	"context"
	"database/sql"
	"encoding/xml"
	"fmt"
	"html"
	"io"
	"net/http"
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

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Items       []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
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

func FetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "gator")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	xmlData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var feed RSSFeed
	if err := xml.Unmarshal(xmlData, &feed); err != nil {
		return nil, err
	}
	html.UnescapeString(feed.Channel.Title)
	html.UnescapeString(feed.Channel.Description)
	for _, item := range feed.Channel.Items {
		html.UnescapeString(item.Title)
		html.UnescapeString(item.Description)
	}
	return &feed, nil
}

func HandlerAgg(s *State, cmd Command) error {
	feedURL := "https://www.wagslane.dev/index.xml"
	feed, err := FetchFeed(context.Background(), feedURL)
	if err != nil {
		return err
	}
	if feed == nil {
		return fmt.Errorf("feed is nil")
	}
	fmt.Printf("Feed Title: %s\n", feed.Channel.Title)
	fmt.Printf("Feed Description: %s\n", feed.Channel.Description)
	for _, item := range feed.Channel.Items {
		fmt.Printf("Item Title: %s\n", item.Title)
		fmt.Printf("Item Description: %s\n", item.Description)
	}
	return nil
}

func HandlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) != 2 {
		fmt.Println("Usage: addfeed <feed_name> <feed_url>")
		os.Exit(1)
	}
	feedName := cmd.Args[0]
	feedURL := cmd.Args[1]
	user, err := s.DB.GetUser(context.Background(), s.Config.CurrentUserName)
	if err != nil {
		return err
	}
	feedParams := database.AddFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      feedName,
		Url:       feedURL,
		UserID:    user.ID,
	}
	feed, err := s.DB.AddFeed(context.Background(), feedParams)
	if err != nil {
		return err
	}
	fmt.Printf("Feed added successfully\n")
	fmt.Printf("Feed ID: %s\n", feed.ID)
	fmt.Printf("Feed CreatedAt: %s\n", feed.CreatedAt)
	fmt.Printf("Feed UpdatedAt: %s\n", feed.UpdatedAt)
	fmt.Printf("Feed Name: %s\n", feed.Name)
	fmt.Printf("Feed URL: %s\n", feed.Url)
	fmt.Printf("Feed User: %s\n", s.Config.CurrentUserName)
	fmt.Printf("Feed UserID: %s\n", feed.UserID)
	return nil
}
