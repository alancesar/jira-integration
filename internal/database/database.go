package database

import (
	"context"
	"gorm.io/gorm"
	"jira-integration/internal/database/model"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/sprint"
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
		&model.Project{},
		&model.StatusCategory{},
		&model.Status{},
		&model.IssueType{},
		&model.Sprint{},
		&model.Product{},
		&model.Changelog{},
		&model.Issue{},
	); err != nil {
		log.Fatal("while running auto migrate", err)
	}

	return &SQLite{
		db: db,
	}
}

func (l SQLite) SaveIssue(ctx context.Context, i issue.Issue) error {
	if exists, err := l.exists(ctx, i); err != nil {
		return err
	} else if exists {
		return l.updateIssue(ctx, i)
	}

	return l.insertIssue(ctx, i)
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

func (l SQLite) SaveIssueType(ctx context.Context, it issue.Type) error {
	m := model.NewIssueType(it)
	tx := l.db.WithContext(ctx).Save(&m)
	return tx.Error
}

func (l SQLite) SaveStatus(ctx context.Context, s issue.Status) error {
	m := model.NewStatus(s)
	tx := l.db.WithContext(ctx).Save(&m)
	return tx.Error
}

func (l SQLite) SaveSprint(ctx context.Context, s sprint.Sprint) error {
	m := model.NewSprint(s)
	tx := l.db.WithContext(ctx).Save(&m)
	return tx.Error
}

func (l SQLite) SaveChangelog(ctx context.Context, i issue.Issue, c issue.Changelog) error {
	m := model.NewChangelog(i, c)
	tx := l.db.WithContext(ctx).Save(&m)
	return tx.Error
}

func (l SQLite) exists(ctx context.Context, i issue.Issue) (bool, error) {
	var exists bool
	tx := l.db.Model(&model.Issue{}).
		WithContext(ctx).
		Select("count(*) > 0").
		Where("id = ?", i.ID).
		Find(&exists)

	return exists, tx.Error
}

func (l SQLite) insertIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	tx := l.db.Omit("Parent").WithContext(ctx).Create(&m)
	return tx.Error
}

func (l SQLite) updateIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	sprints := m.Sprints
	if err := l.db.Model(&m).Association("Sprints").Clear(); err != nil {
		return err
	}

	products := m.Products
	if err := l.db.Model(&m).Association("Products").Clear(); err != nil {
		return err
	}

	m.Sprints = sprints
	m.Products = products
	tx := l.db.Omit("Parent").WithContext(ctx).Save(&m)
	return tx.Error
}
