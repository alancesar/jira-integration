package usecase

import (
	"context"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"jira-integration/internal/database"
	"jira-integration/internal/jira"
	"jira-integration/pkg/issue"
	"net/http"
	"testing"
)

func TestIssueUseCase_GetIssueByKey(t *testing.T) {
	type args struct {
		ctx context.Context
		key string
	}
	tests := []struct {
		name    string
		fields  IssueClient
		args    args
		want    issue.Issue
		wantErr bool
	}{
		{
			name: "some test",
			args: args{
				ctx: context.Background(),
				key: "DFX-4980",
			},
			wantErr: false,
		},
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatal(err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			uc := FetchUseCase{
				client: jira.NewClient(0, jira.Credentials{
					Username: "alan.elias@ebury.com",
					Password: "ATATT3xFfGF03x1ofn53rD2PwPzw_-ohqX1Xu1WzMTQ-qFFW7F-kFuaKOfPq2KfyToAvgBG0GE04O5sYz2ybUoFOKmMt6F-a9tn52qowJecChjF6jRmYm-jKH65gLgaQcTJcXRxAR0UKB-Fqp11GgvQ9nBqK6nzRTNGnqY-CkW6WIRYCUuDjYWc=29C7F8F7",
				}, http.DefaultClient),
				db: database.NewSQLite(db),
			}
			if err := uc.Execute(tt.args.ctx, tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
