package database

import (
	"context"
	"errors"
	"gorm.io/gorm"
	"jira-integration/internal/database/model"
	"jira-integration/pkg/issue"
	"log"
)

type (
	Gorm struct {
		db *gorm.DB
	}
)

func NewGorm(db *gorm.DB) *Gorm {
	if err := db.AutoMigrate(
		&model.Issue{},
		&model.Changelog{},
		&model.Sprint{},
	); err != nil {
		log.Fatalln("while running auto migrate", err)
	}

	return &Gorm{
		db: db,
	}
}

func (l Gorm) CreateIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	if err := l.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (l Gorm) UpdateIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	if err := l.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}

	return nil
}

func (l Gorm) IssueExistsByKey(ctx context.Context, issueKey string) (bool, error) {
	err := l.db.WithContext(ctx).
		Where("key = ?", issueKey).
		First(&model.Issue{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}
