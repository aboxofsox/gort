package gort

import "sync"

type Store struct {
	mu    sync.Mutex
	Items map[string]any
}

func NewStore() *Store {
	return &Store{
		Items: make(map[string]any),
	}
}

func (s *Store) Set(key string, value any) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Items[key] = value
}

func (s *Store) Get(key string) (any, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.Items[key]; !ok {
		return nil, false
	}
	return s.Items[key], true
}

func (s *Store) Remove(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.Items, key)
}

func (s *Store) Purge() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.Items = make(map[string]any)
}

func (s *Store) Size() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.Items)
}
