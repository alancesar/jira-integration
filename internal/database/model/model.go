package model

import (
	"encoding/json"
	"gorm.io/datatypes"
	"jira-integration/pkg/financial"
	"jira-integration/pkg/issue"
	"time"
)

type (
	Changelog struct {
		ID        uint `gorm:"primarykey"`
		IssueID   uint
		From      string
		To        string
		CreatedAt time.Time `gorm:"autoCreateTime:false"`
	}

	Issue struct {
		ID          uint   `gorm:"primarykey"`
		Key         string `gorm:"unique"`
		Summary     string
		Status      string
		IssueType   string
		Project     string
		Parent      *string
		Sprint      *string
		Labels      datatypes.JSON
		Assignee    *string
		Reporter    string
		StoryPoints *uint
		Products    datatypes.JSON
		FixVersion  *string
		Locality    *string
		Changelog   []Changelog
		CreatedAt   time.Time `gorm:"autoCreateTime:false"`
		UpdatedAt   time.Time `gorm:"autoUpdateTime:false"`
	}

	Volume struct {
		PartnerID   string `gorm:"primaryKey"`
		CreatedAt   string `gorm:"primaryKey"`
		Offer       string `gorm:"primaryKey"`
		Product     string `gorm:"primaryKey"`
		PartnerName string
		Volume      float64
		Operations  int
	}
)

func NewIssue(i issue.Issue) *Issue {
	labels, _ := json.Marshal(i.Labels)
	products, _ := json.Marshal(i.Products)

	changelog := make([]Changelog, len(i.Changelog), len(i.Changelog))
	for index, c := range i.Changelog {
		changelog[index] = NewChangelog(c, i.ID)
	}

	return &Issue{
		ID:          i.ID,
		Key:         i.Key,
		Summary:     i.Summary,
		Status:      i.Status,
		IssueType:   i.IssueType,
		Project:     i.Project,
		Parent:      stringToPointer(i.Parent),
		Sprint:      stringToPointer(i.Sprint),
		Labels:      labels,
		Assignee:    stringToPointer(i.Assignee),
		Reporter:    i.Reporter,
		StoryPoints: uintToPointer(i.StoryPoints),
		Products:    products,
		FixVersion:  stringToPointer(i.FixVersion),
		Locality:    stringToPointer(i.Locality),
		Changelog:   changelog,
		CreatedAt:   i.CreatedAt,
		UpdatedAt:   i.UpdatedAt,
	}
}

func NewChangelog(c issue.Changelog, issueID uint) Changelog {
	return Changelog{
		ID:        c.ID,
		IssueID:   issueID,
		From:      c.From,
		To:        c.To,
		CreatedAt: c.CreatedAt,
	}
}

func NewVolume(v financial.Volume) Volume {
	return Volume{
		PartnerID:   v.PartnerID,
		CreatedAt:   v.CreatedAt,
		Offer:       v.Offer,
		Product:     v.Product,
		PartnerName: v.PartnerName,
		Volume:      v.Volume,
		Operations:  v.Operations,
	}
}

func stringToPointer(value string) *string {
	if value == "" {
		return nil
	}

	return &value
}

func uintToPointer(value uint) *uint {
	if value == 0 {
		return nil
	}

	return &value
}
