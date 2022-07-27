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
	GithubId int64
	Name string
	Owner string
	FullName string
	URL string
	CreatedAt time.Time
	Topics []string
	Languages map[string]int
}

RepositoryStats struct {
	Language map[string]float64
	Owner map[string]float64
	Topic map[string]float64
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