package main

import (
	"context"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jira-integration/internal/database"
	"jira-integration/internal/jira"
	"jira-integration/usecase"
	"log"
	"net/http"
	"os"
)

func init() {
	if len(os.Args) < 2 {
		log.Fatalln("missing status argument")
	}
}

func main() {
	ctx := context.Background()
	dsn := os.Getenv("JIRA_DB_DSN")
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}

	jiraClient := jira.NewClient(jira.Credentials{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}, http.DefaultClient)

	postgresDB := database.NewGorm(conn)

	syncSprintsUseCase := usecase.NewSyncSprintsUseCase(jiraClient, postgresDB)
	if err := syncSprintsUseCase.Execute(ctx, os.Args[1:]); err != nil {
		log.Fatalln(err)
	}
}
