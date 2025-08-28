package main

import (
	"context"
	"flag"
	"fmt"
	"jira-integration/internal/database"
	"jira-integration/internal/jira"
	"jira-integration/usecase"
	"log"
	"net/http"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	jql string
)

func init() {
	flag.StringVar(&jql, "jql", "", "JQL query")
	flag.Parse()

	if jql == "" {
		log.Fatalln("jql query is required")
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
	fetchUseCase := usecase.NewFetchUseCase(jiraClient, postgresDB)
	streamUseCase := usecase.NewStreamUseCase(jiraClient, fetchUseCase.Execute, postgresDB)
	fmt.Println("fetching issues with JQL:", jql)
	if err := streamUseCase.Execute(ctx, jql); err != nil {
		log.Fatalln(err)
	}
}
