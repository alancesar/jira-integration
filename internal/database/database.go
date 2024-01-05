package database

import (
	"context"
	"gorm.io/gorm"
	"jira-integration/internal/database/model"
	"jira-integration/pkg/issue"
	"jira-integration/pkg/sprint"
	"time"
)

type (
	SQLite struct {
		db *gorm.DB
	}
)

func NewSQLite(db *gorm.DB) *SQLite {
	_ = db.AutoMigrate(
		&model.Project{},
		&model.StatusCategory{},
		&model.Status{},
		&model.Resolution{},
		&model.FixVersion{},
		&model.Priority{},
		&model.IssueType{},
		&model.Sprint{},
		&model.Issue{},
	)
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
	var updatedAt time.Time
	tx := l.db.WithContext(ctx).Raw(`select max(updated_at) from issues`).Scan(&updatedAt)
	if tx.Error != nil {
		return time.Time{}, tx.Error
	}

	return updatedAt, nil
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

func (l SQLite) SaveFixVersion(ctx context.Context, fv issue.FixVersion) error {
	m := model.NewFixVersion(fv)
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

	fixVersions := m.FixVersions
	if err := l.db.Model(&m).Association("FixVersions").Clear(); err != nil {
		return err
	}

	m.Sprints = sprints
	m.FixVersions = fixVersions
	tx := l.db.Omit("Parent").WithContext(ctx).Save(&m)
	return tx.Error
}
