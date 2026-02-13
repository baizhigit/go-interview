package main

import (
	"fmt"
	"sort"
	"sync"
)

type Document struct {
	DocID   string
	Content string
}

type User struct {
	UserID string
}

type Scorer interface {
	GetScore(doc Document, user User) int
}

type DocumentService interface {
	AddDocument(doc Document) error
	GetTopDocuments(user User, limit int) ([]Document, error)
}

type DocumentServiceImpl struct {
	docs      map[string]Document
	cache     map[string][]Document
	cacheSize int
	mu        sync.RWMutex
	scorer    Scorer
}

func NewDocumentService(scorer Scorer, cacheSize int) *DocumentServiceImpl {
	return &DocumentServiceImpl{
		docs:      make(map[string]Document),
		cache:     make(map[string][]Document),
		cacheSize: cacheSize,
		scorer:    scorer,
	}
}

func (s *DocumentServiceImpl) AddDocument(doc Document) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.docs[doc.DocID] = doc
	s.cache = make(map[string][]Document)
	return nil
}

func (s *DocumentServiceImpl) GetTopDocuments(user User, limit int) ([]Document, error) {
	s.mu.RLock()

	if cached, ok := s.cache[user.UserID]; ok {
		s.mu.RUnlock()
		if limit > len(cached) {
			return cached, nil
		}
		return cached[:limit], nil
	}
	s.mu.RUnlock()

	s.mu.Lock()
	defer s.mu.Unlock()

	if cached, ok := s.cache[user.UserID]; ok {
		s.mu.RUnlock()
		if limit > len(cached) {
			return cached, nil
		}
		return cached[:limit], nil
	}

	type scoredDoc struct {
		doc   Document
		score int
	}

	scores := make([]scoredDoc, len(s.docs))
	for _, doc := range s.docs {
		score := s.scorer.GetScore(doc, user)
		scores = append(scores, scoredDoc{
			doc:   doc,
			score: score,
		})
	}

	sort.Slice(scores, func(i, j int) bool {
		return scores[i].score > scores[j].score
	})

	res := make([]Document, 0, s.cacheSize)
	for i := 0; i < s.cacheSize && i < len(scores); i++ {
		res = append(res, scores[i].doc)
	}
	s.cache[user.UserID] = res

	if limit > len(res) {
		return res, nil
	}
	return res[:limit], nil
}

func main() {
	fmt.Println("main start")

	fmt.Println("main end")
}
