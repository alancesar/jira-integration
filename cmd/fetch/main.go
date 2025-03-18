package main

import (
	"context"
	"flag"
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

var (
	jql string
)

func init() {
	flag.StringVar(&jql, "jql", "", "JQL query")
	flag.Parse()
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
	fetchUseCase := usecase.NewFetchUseCase(jiraClient, postgresDB)
	streamUseCase := usecase.NewStreamUseCase(jiraClient, fetchUseCase.Execute)
	if err := streamUseCase.Execute(ctx, jql); err != nil {
		log.Fatalln(err)
	}
}
