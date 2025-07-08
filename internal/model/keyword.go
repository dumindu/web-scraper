package model

import (
	"github.com/google/uuid"
)

type Keywords []*Keyword

type Keyword struct {
	Model2
	UserID       uuid.UUID
	Keyword      string
	Status       string
	SearchEngine string
	AdCount      *int64
	LinkCount    *int64
	HTMLContent  *string
	ErrorMessage *string
}

type KeywordDTO struct {
	ID           int64   `json:"id"`
	Keyword      string  `json:"keyword"`
	Status       string  `json:"status"`
	AdCount      *int64  `json:"adCount"`
	LinkCount    *int64  `json:"linkCount"`
	HTMLContent  *string `json:"htmlContent"`
	ErrorMessage *string `json:"errorMessage"`
}

func (ks Keywords) ToDTOs() []*KeywordDTO {
	result := make([]*KeywordDTO, len(ks))
	for i, v := range ks {
		result[i] = v.ToDTO()
	}

	return result
}

func (k *Keyword) ToDTO() *KeywordDTO {
	return &KeywordDTO{
		ID:           k.ID,
		Keyword:      k.Keyword,
		Status:       k.Status,
		AdCount:      k.AdCount,
		LinkCount:    k.LinkCount,
		HTMLContent:  k.HTMLContent,
		ErrorMessage: k.ErrorMessage,
	}
}
