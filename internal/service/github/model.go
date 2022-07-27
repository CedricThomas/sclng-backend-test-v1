package github

import (
	"encoding/json"
	"time"
)

type (

ListingFilters struct {
	Language []string
	Owner []string
	Topic []string
}

Repository struct {
	GithubId int64 `json:"github_id"`
	Name string `json:"name"`
	Owner string `json:"owner"`
	FullName string `json:"fullName"`
	URL string `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	Topics []string `json:"topics"`
	Languages map[string]int `json:"languages"`
}

RepositoryStats struct {
	Languages map[string]float64 `json:"languages"`
	Owners map[string]float64 `json:"owners"`
	Topics map[string]float64 `json:"topics"`
}

)

func NewListingFilters(rawFilters map[string][]string) (*ListingFilters, error) {
	var filters ListingFilters
	dbByte, err := json.Marshal(rawFilters)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(dbByte, &filters)
	return &filters, err
}