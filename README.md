# 🐊 Gator - RSS Feed Aggregator

**Gator** is a powerful command-line RSS feed aggregator built with **Go** and **PostgreSQL**. It provides a complete user management system, RSS feed tracking, and post aggregation with type-safe database operations using SQLC.

## ✨ Features

- **User Management**: Register, login, and manage multiple users
- **Feed Management**: Add, follow, and unfollow RSS feeds
- **Post Aggregation**: Automatically fetch and store posts from followed feeds
- **Browse Posts**: View recent posts with customizable limits
- **Type-Safe Database**: SQLC-generated Go code for safe and efficient database operations
- **PostgreSQL Backend**: Robust relational database with proper migrations
- **JSON Configuration**: Simple home directory configuration management

## 🏗️ Project Structure

```
gator/
├── internal/
│   ├── commands/         # CLI command handlers and RSS parsing
│   ├── config/           # JSON configuration management
│   └── database/         # SQLC-generated Go database code
├── sql/
│   ├── queries/          # SQL queries (users.sql, feeds.sql, posts.sql)
│   └── schema/           # Database migration files (001-005)
├── main.go              # Application entry point
├── sqlc.yaml            # SQLC configuration
├── go.mod               # Go module definition
└── README.md            # This file
```

## 🚀 Getting Started

### Prerequisites

- Go 1.24.4 or higher
- PostgreSQL database
- [SQLC](https://sqlc.dev/) for code generation
- [Goose](https://github.com/pressly/goose) for database migrations (optional)

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/UUest/gator.git
   cd gator
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Set up your database**
   - Create a PostgreSQL database
   - Run migrations from `sql/schema/` directory
   
4. **Configure the application**
   
   Create a `.gatorconfig.json` file in your home directory:
   ```json
   {
     "db_url": "postgres://username:password@localhost:5432/gator?sslmode=disable",
     "current_user_name": ""
   }
   ```

5. **Generate database code**
   ```bash
   sqlc generate
   ```

6. **Build the application**
   ```bash
   go build -o gator
   ```

## 📖 Usage

### User Management

```bash
# Register a new user
./gator register <username>

# Login as existing user
./gator login <username>

# List all users (* indicates current user)
./gator users

# Reset database (removes all data)
./gator reset
```

### Feed Management

```bash
# Add a new RSS feed (automatically follows it)
./gator addfeed "<feed_name>" "<feed_url>"

# List all feeds
./gator feeds

# Follow an existing feed
./gator follow "<feed_url>"

# List feeds you're following
./gator following

# Unfollow a feed
./gator unfollow "<feed_url>"
```

### Content Aggregation

```bash
# Start continuous feed aggregation (fetches every interval)
./gator agg <duration>  # e.g., "30s", "5m", "1h"

# Browse recent posts (default: 2 posts)
./gator browse

# Browse specific number of posts
./gator browse <limit>
```

### Example Workflow

```bash
# 1. Register and login
./gator register alice
./gator login alice

# 2. Add some feeds
./gator addfeed "TechCrunch" "https://techcrunch.com/feed/"
./gator addfeed "Hacker News" "https://hnrss.org/frontpage"

# 3. Start aggregation in background
./gator agg 10m &

# 4. Browse posts
./gator browse 5
```

## 🗄️ Database Schema

The application uses the following database tables:

- **users**: User accounts with UUID primary keys
- **feeds**: RSS feed metadata and ownership
- **feed_follows**: Many-to-many relationship between users and feeds
- **posts**: Aggregated RSS feed items with publication dates
- **last_fetched**: Tracking when feeds were last updated

## 🔧 Configuration

### Database URL Format
```
postgres://username:password@host:port/database?sslmode=disable
```

### Configuration File Location
- **Linux/macOS**: `~/.gatorconfig.json`
- **Windows**: `%USERPROFILE%\.gatorconfig.json`

### Environment Variables
You can also set the database URL via environment variable:
```bash
export DATABASE_URL="postgres://username:password@localhost:5432/gator"
```

## 🛠️ Development

### Database Migrations

Using [Goose](https://github.com/pressly/goose):

```bash
# Install goose
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir sql/schema postgres "$DATABASE_URL" up

# Check migration status
goose -dir sql/schema postgres "$DATABASE_URL" status
```

### Code Generation

```bash
# Regenerate database code after schema changes
sqlc generate
```

### Dependencies

- **Core**: Go standard library
- **Database**: `github.com/lib/pq` (PostgreSQL driver)
- **UUID**: `github.com/google/uuid`
- **Code Generation**: [SQLC](https://sqlc.dev/)

## 🎯 Roadmap

### Immediate Improvements
- [ ] **Help Command**: Add `gator help` to show all available commands
- [ ] **Better Error Handling**: More descriptive error messages and recovery
- [ ] **Configuration Validation**: Validate database connection on startup
- [ ] **Logging**: Add structured logging for debugging and monitoring

### Feature Enhancements
- [ ] **OPML Import/Export**: Import/export feed lists in OPML format
- [ ] **Feed Categories**: Organize feeds into categories/tags
- [ ] **Search Functionality**: Search through aggregated posts
- [ ] **Duplicate Detection**: Prevent duplicate posts from being stored
- [ ] **Feed Health Monitoring**: Track feed availability and errors

### Advanced Features
- [ ] **Web Interface**: Optional HTTP server for browser-based interaction
- [ ] **Webhooks**: Notify external services when new posts arrive
- [ ] **Feed Discovery**: Auto-discover RSS feeds from website URLs
- [ ] **Content Filtering**: Filter posts by keywords or patterns
- [ ] **Export Formats**: Export posts to JSON, CSV, or other formats

### Performance & Reliability
- [ ] **Concurrent Fetching**: Fetch multiple feeds simultaneously
- [ ] **Retry Logic**: Implement exponential backoff for failed requests
- [ ] **Rate Limiting**: Respect feed rate limits and etiquette
- [ ] **Caching**: Cache feed content to reduce bandwidth usage
- [ ] **Database Indexing**: Optimize database queries for large datasets

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by RSS aggregators and feed readers
- Built with [SQLC](https://sqlc.dev/) for type-safe database operations
- Database migrations powered by [Goose](https://github.com/pressly/goose)
- PostgreSQL driver by [lib/pq](https://github.com/lib/pq)