package github

import (
	"context"
	"sync"

	"github.com/google/go-github/v45/github"
	"golang.org/x/oauth2"
	"golang.org/x/sync/errgroup"
)

type Service struct {
	client *github.Client
}

func NewService(GithubApiKey string) (*Service, error) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: GithubApiKey},
	)
	tc := oauth2.NewClient(ctx, ts)
	return &Service{
		github.NewClient(tc),
	}, nil
}

func (s *Service) ListPublicRepositories(ctx context.Context, filters *ListingFilters) ([]*Repository, error) {
	errs, ctx := errgroup.WithContext(ctx)
	rawRepos, _, err := s.client.Repositories.ListAll(ctx, nil)
	mx := sync.Mutex{}
	if err != nil {
		return nil, err
	}
	repos := make([]*Repository, 0, len(rawRepos))
	for _, repo := range rawRepos {
		repo := repo
		errs.Go(func() error {
			languageStats, _, err := s.client.Repositories.ListLanguages(ctx, repo.GetOwner().GetLogin(), repo.GetName())
			if err != nil {
				return err
			}
			topicsStats, _, err := s.client.Repositories.ListAllTopics(ctx, repo.GetOwner().GetLogin(), repo.GetName())
			if err != nil {
				return err
			}
			mx.Lock()
			repos = append(repos, &Repository{
				GithubId:     repo.GetID(),
				Name:         repo.GetName(),
				Owner:        repo.GetOwner().GetLogin(),
				FullName:     repo.GetFullName(),
				URL:          repo.GetHTMLURL(),
				CreatedAt:	  repo.GetCreatedAt().Time,
				Topics:       topicsStats,
				Languages:    languageStats,
			})
			mx.Unlock()
			return nil
		})
	}
	err = errs.Wait()
	if err != nil {
		return nil, err
	}
	return filterRepositories(repos, filters), err
}

func (s *Service) GetPublicRepositoriesStats(ctx context.Context, filter *ListingFilters) (*RepositoryStats, error) {
	stats := RepositoryStats{
		Language: map[string]float64{},
		Owner:    map[string]float64{},
		Topic:    map[string]float64{},
	}
	repos, err := s.ListPublicRepositories(ctx, filter)
	if err != nil {
		return nil, err
	}
	nbRepos := float64(len(repos))
	var totalBytes float64
	pushToKey := func(m map[string]float64, k string, v float64) {
		if val, ok := m[k]; ok {
			m[k] = val + v
		} else {
			m[k] = v
		}
	}
	divideMap := func(m map[string]float64, max float64) {
		for key, value := range m {
			m[key] = roundFloat(value / max, 2)
		}
	}
	for _, repo := range repos {
		pushToKey(stats.Owner, repo.Owner, 1.0)
		for _, topic := range repo.Topics {
			pushToKey(stats.Topic, topic, 1.0)
		}
		for language, value := range repo.Languages {
			totalBytes += float64(value)
			pushToKey(stats.Language, language, float64(value))
		}
	}
	divideMap(stats.Owner, nbRepos)
	divideMap(stats.Topic, nbRepos)
	divideMap(stats.Language, totalBytes)
	return &stats, nil
}








