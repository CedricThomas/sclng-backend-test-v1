package github

import "math"

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func filterRepositories(repos []*Repository, filter *ListingFilters) []*Repository {
	f := func(repos []*Repository, dict []string, test func(*Repository, string) bool) (ret []*Repository) {
		if len(dict) == 0 {
			ret = repos
			return
		}
		for _, repo := range repos {
			for _, dictValue := range dict {
				if test(repo, dictValue) {
					ret = append(ret, repo)
					break
				}
			}
		}
		return
	}
	repos = f(repos, filter.Owner, func(repo *Repository, value string) bool {
		return repo.Owner == value
	})
	repos = f(repos, filter.Language, func(repo *Repository, value string) bool {
		_, ok := repo.Languages[value]
		return ok
	})
	repos = f(repos, filter.Topic, func(repo *Repository, value string) bool {
		for _, topic := range repo.Topics {
			if topic == value {
				return true
			}
		}
		return false
	})
	return repos
}
