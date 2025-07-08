package tasks

import (
	"fmt"

	"github.com/hibiken/asynq"
)

const (
	TypeScrapeKeyword           = "scrape:keyword"
	fmtScrapeKeywordPayloadJSON = `{"keywordID": %d}`
	ScrapeKeywordDelayInSeconds = 1 // Enqueue with a delay to avoid rate limiting
)

type ScrapeKeywordPayload struct {
	KeywordID int64 `json:"keywordID"`
}

func NewScrapeKeywordTask(keywordID int64) *asynq.Task {
	payload := []byte(fmt.Sprintf(fmtScrapeKeywordPayloadJSON, keywordID))
	return asynq.NewTask(TypeScrapeKeyword, payload)
}
