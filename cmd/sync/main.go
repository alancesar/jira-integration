package main

import (
	"context"
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"jira-integration/gateway"
	"jira-integration/internal/database"
	"jira-integration/pkg/jira"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	full    bool
	query   string
	boardID int
)

func init() {
	flag.BoolVar(&full, "full", false, "run sync for all issues")
	flag.StringVar(&query, "project", "project = DFX", "project key")
	flag.IntVar(&boardID, "board-id", 66, "board id")
	flag.Parse()
}

func main() {
	ctx := context.Background()

	gormConnection, err := connectionFactory()
	if err != nil {
		log.Fatalln(err)
	}

	db := database.NewSQLite(gormConnection)
	args := []string{query}
	if !full {
		days := retrieveDaysSinceLastUpdate(ctx, db)
		args = append(args, fmt.Sprintf("updated >= -%dd", days))
	}

	credentials := jira.Credentials{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}

	client := jira.NewClient(boardID, credentials, http.DefaultClient)
	g := gateway.New(client)

	if full {
		if err := g.SyncDependencies(ctx); err != nil {
			log.Fatalln(err)
		}
	}

	g.SyncIssues(ctx, args...)
}

func retrieveDaysSinceLastUpdate(ctx context.Context, db *database.SQLite) int {
	lastUpdate, err := db.GetLastUpdate(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return int(time.Since(lastUpdate).Hours()/24) + 1
}

func connectionFactory() (*gorm.DB, error) {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn == "" {
		return gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{})
	} else {
		return gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}
}
