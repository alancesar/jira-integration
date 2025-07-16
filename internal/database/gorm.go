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
		&model.Label{},
		&model.Product{},
		&model.Changelog{},
		&model.Sprint{},
		&model.Account{},
		&model.Issue{},
	); err != nil {
		log.Fatalln("while running auto migrate", err)
	}

	return &Gorm{
		db: db,
	}
}

func (g Gorm) CreateIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	if err := g.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}

	return nil
}

func (g Gorm) UpdateIssue(ctx context.Context, i issue.Issue) error {
	m := model.NewIssue(i)
	if err := g.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}

	return nil
}

func (g Gorm) IssueExists(ctx context.Context, i issue.Issue) (bool, error) {
	err := g.db.WithContext(ctx).
		Where("id = ? OR key = ?", i.ID, i.Key).
		First(&model.Issue{}).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func (g Gorm) GetSprintsByState(ctx context.Context, states []string) ([]issue.Sprint, error) {
	var sprints []model.Sprint
	if err := g.db.WithContext(ctx).Where("state in (?)", states).Find(&sprints).Error; err != nil {
		return nil, err
	}

	return model.Sprints(sprints).ToDomain(), nil
}

func (g Gorm) SaveSprint(ctx context.Context, sprint issue.Sprint) error {
	m := model.NewSprint(&sprint)
	if err := g.db.WithContext(ctx).Save(m).Error; err != nil {
		return err
	}

	return nil
}
