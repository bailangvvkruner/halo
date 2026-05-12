package service

import (
	"context"
	"sync"
)

type StatsService interface {
	IncrVisit(ctx context.Context, name string) error
	GetVisitCount(ctx context.Context, name string) (int64, error)
	GetAllStats(ctx context.Context) (map[string]int64, error)
}

type statsService struct {
	mu     sync.RWMutex
	counts map[string]int64
}

func NewStatsService() StatsService {
	return &statsService{counts: make(map[string]int64)}
}

func (s *statsService) IncrVisit(_ context.Context, name string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.counts[name]++
	return nil
}

func (s *statsService) GetVisitCount(_ context.Context, name string) (int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	count, ok := s.counts[name]
	if !ok {
		return 0, nil
	}
	return count, nil
}

func (s *statsService) GetAllStats(_ context.Context) (map[string]int64, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make(map[string]int64, len(s.counts))
	for k, v := range s.counts {
		result[k] = v
	}
	return result, nil
}
