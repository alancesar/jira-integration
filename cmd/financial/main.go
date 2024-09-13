package main

import (
	"context"
	"encoding/csv"
	"flag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jira-integration/internal/database"
	"jira-integration/pkg/financial"
	"log"
	"os"
)

var (
	source string
)

func init() {
	flag.StringVar(&source, "source", "report.csv", "source csv file")
	flag.Parse()
}

func main() {
	ctx := context.Background()
	gormConnection, err := gorm.Open(sqlite.Open("sqlite.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalln(err)
	}

	db := database.NewSQLite(gormConnection)
	reportFile, err := os.Open("./report.csv")
	if err != nil {
		log.Fatalln(err)
	}

	reader := csv.NewReader(reportFile)
	records, err := reader.ReadAll()
	if err != nil {
		log.Fatalln(err)
	}

	for i, record := range records {
		if i == 0 {
			continue
		}

		v, err := financial.NewVolume(record)
		if err != nil {
			log.Println(err)
		}

		if err := db.SaveFinancial(ctx, v); err != nil {
			log.Println(err)
		}
	}
}
