package database

import (
	"context"
	"gorm.io/gorm"
	"jira-integration/internal/database/model"
	"jira-integration/pkg/financial"
	"jira-integration/pkg/issue"
	"log"
	"time"
)

type (
	SQLite struct {
		db *gorm.DB
	}
)

func NewSQLite(db *gorm.DB) *SQLite {
	if err := db.AutoMigrate(
		&model.Changelog{},
		&model.Issue{},
		&model.Volume{},
	); err != nil {
		log.Fatalln("while running auto migrate", err)
	}

	return &SQLite{
		db: db,
	}
}

func (l SQLite) SaveIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	if err := l.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}

	return nil
}

func (l SQLite) GetLastUpdate(ctx context.Context) (time.Time, error) {
	var rawDatetime string
	tx := l.db.WithContext(ctx).Raw(`select max(updated_at) from issues`).Scan(&rawDatetime)
	if tx.Error != nil {
		return time.Time{}, tx.Error
	} else if rawDatetime == "" {
		return time.Time{}, nil
	}

	rawDate := rawDatetime[0:10]
	parsed, err := time.Parse("2006-01-02", rawDate)
	if err != nil {
		return time.Time{}, err
	}

	return parsed, nil
}

func (l SQLite) SaveFinancial(ctx context.Context, v financial.Volume) error {
	m := model.NewVolume(v)
	tx := l.db.WithContext(ctx).Save(&m)
	return tx.Error
}

func (l SQLite) issueExists(ctx context.Context, i issue.Issue) (bool, error) {
	var exists bool
	tx := l.db.Model(&model.Issue{}).
		WithContext(ctx).
		Select("count(*) > 0").
		Where("id = ?", i.ID).
		Find(&exists)

	return exists, tx.Error
}
