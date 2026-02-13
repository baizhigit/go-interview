// Задача:
// Реализовать сервис для сбора аналитики в реальном времени.
// Сервис получает события от пользователей и должен уметь отдавать
// количество событий за последние 5 минут для каждого пользователя.
//
// Требования:
// - HandleEvent вызывается при каждом событии пользователя
// - GetCount возвращает количество событий за последние 5 минут
// - Учитывать конкурентный доступ (метод может вызываться из горутин)
// - Минимизировать потребление памяти (удалять устаревшие события)
// - Оптимизировать производительность GetCount

package main

import (
	"fmt"
	"sync"
	"time"
)

type Service interface {
	HandleEvent(userName string, currentTime time.Time)
	GetCount(userName string, currentTime time.Time) int
}

type analyticsService struct {
	mu     sync.Mutex
	events map[string][]time.Time
	window time.Duration
}

func NewAnalyticsService(window int) Service {
	return &analyticsService{
		events: make(map[string][]time.Time),
		window: time.Duration(window) * time.Minute,
	}
}

func (s *analyticsService) HandleEvent(userName string, currentTime time.Time) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[userName] = append(s.events[userName], currentTime)
}

func (s *analyticsService) GetCount(userName string, currentTime time.Time) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	timestamps, exists := s.events[userName]
	if !exists {
		return 0 // Пользователь не найден
	}

	// Вычисляем границу временного окна.
	windowStart := currentTime.Add(-s.window)

	// Находим первый элемент, который попадает в окно.
	// Все элементы до него можно удалить.
	validStart := 0
	for i, ts := range timestamps {
		if ts.After(windowStart) || ts.Equal(windowStart) {
			validStart = i
			break
		}
	}

	// Удаляем устаревшие события для экономии памяти.
	// Создаем новый слайс с валидными событиями.
	validEvents := timestamps[validStart:]
	s.events[userName] = validEvents

	// Возвращаем количество событий в окне.
	return len(validEvents)
}

func main() {
	service := NewAnalyticsService(5)

	now := time.Now()

	// Регистрируем события
	service.HandleEvent("user1", now.Add(-6*time.Minute)) // Старое, не попадет
	service.HandleEvent("user1", now.Add(-4*time.Minute))
	service.HandleEvent("user1", now.Add(-3*time.Minute))
	service.HandleEvent("user1", now.Add(-1*time.Minute))
	service.HandleEvent("user2", now.Add(-2*time.Minute))

	// Получаем статистику
	fmt.Printf("user1 events: %d\n", service.GetCount("user1", now)) // 3
	fmt.Printf("user2 events: %d\n", service.GetCount("user2", now)) // 1
	fmt.Printf("user3 events: %d\n", service.GetCount("user3", now)) // 0
}
