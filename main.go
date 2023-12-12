package main

import (
	"context"
	"flag"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"jira-integration/internal/database"
	"jira-integration/pkg/gateway"
	"jira-integration/pkg/jira"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	full    bool
	project string
)

func init() {
	flag.BoolVar(&full, "full", false, "run sync for all issues")
	flag.StringVar(&project, "project", "MAQ", "project key")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	sqliteConnection, err := gorm.Open(postgres.Open("host=localhost user=postgres password=Pa55w0rd dbname=postgres port=5432 sslmode=disable TimeZone=America/Sao_Paulo"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}

	db := database.NewSQLite(sqliteConnection)
	args := []string{fmt.Sprintf("project = %s", project)}
	//if !full {
	//	days := retrieveDaysSinceLastUpdate(ctx, db)
	//	args = append(args, fmt.Sprintf("updated >= -%dd", days))
	//}

	credentials := jira.Credentials{
		Username: os.Getenv("JIRA_USERNAME"),
		Password: os.Getenv("JIRA_PASSWORD"),
	}

	client := jira.NewClient(credentials, http.DefaultClient)

	//issueTypes, err := client.GetIssueTypes()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//for _, issueType := range issueTypes {
	//	fmt.Println("fetching issue type", issueType.Name)
	//	if err := db.SaveIssueType(ctx, issueType); err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	//
	//statuses, err := client.GetStatuses()
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//for _, status := range statuses {
	//	fmt.Println("fetching status", status.Name)
	//	if err := db.SaveStatus(ctx, status); err != nil {
	//		log.Fatalln(err)
	//	}
	//}
	//
	boardID := 66
	sprints := client.StreamSprints(boardID)
	for sprint := range sprints {
		fmt.Println("fetching sprint", sprint.Name)
		if err := db.SaveSprint(ctx, sprint); err != nil {
			log.Fatalln(err)
		}
	}
	//
	//fixVersions := client.StreamFixVersions(boardID)
	//for fixVersion := range fixVersions {
	//	fmt.Println("fetching fix version", fixVersion.Name)
	//	if err := db.SaveFixVersion(ctx, fixVersion); err != nil {
	//		log.Fatalln(err)
	//	}
	//}

	g := gateway.New(client)
	issues := g.StreamAllIssues(args...)
	for issue := range issues {
		log.Println("fetching issue", issue.Key)
		if err := db.SaveIssue(ctx, issue); err != nil {
			log.Println(err)
		}
	}
}

func retrieveDaysSinceLastUpdate(ctx context.Context, db *database.SQLite) int {
	lastUpdate, err := db.GetLastUpdate(ctx)
	if err != nil {
		log.Fatalln(err)
	}

	return int(time.Since(lastUpdate).Hours()/24) + 1
}
