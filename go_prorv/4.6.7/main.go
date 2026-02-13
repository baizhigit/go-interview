package main

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// Этап 1: Последовательные запросы
// ============================================================================

func step1Sequential() {
	fmt.Println("=== Этап 1: Последовательные запросы ===")

	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://eeeboy",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	// Проходим по каждому URL последовательно.
	for _, url := range urls {
		// Создаем GET запрос с таймаутом 5 секунд.
		resp, err := client.Get(url)

		// Проверяем результат запроса.
		if err != nil {
			// Ошибка при выполнении запроса (сеть, DNS и т.д.)
			fmt.Printf("%s - not ok (error: %v)\n", url, err)
			continue
		}

		// Не забываем закрыть тело ответа.
		resp.Body.Close()

		// Проверяем код ответа.
		if resp.StatusCode == http.StatusOK {
			fmt.Printf("%s - ok\n", url)
		} else {
			fmt.Printf("%s - not ok (status: %d)\n", url, resp.StatusCode)
		}
	}
}

// ============================================================================
// Этап 2: Конкурентные запросы с каналами
// ============================================================================

// Result хранит результат проверки URL.
type Result struct {
	URL string
	OK  bool
	Err error
}

func step2Concurrent() {
	fmt.Println("\n=== Этап 2: Конкурентные запросы с каналами ===")

	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://eeeboy",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	// Создаем канал для получения результатов.
	// Буферизированный канал размером len(urls) позволяет горутинам
	// не блокироваться при отправке результатов.
	results := make(chan Result, len(urls))

	// Запускаем горутину для каждого URL.
	for _, url := range urls {
		// ВАЖНО: захватываем url в локальную переменную для горутины.
		// Иначе все горутины будут использовать последнее значение url.
		// Актуально для версий Go до 1.22
		url := url

		go func() {
			// Выполняем запрос в горутине.
			resp, err := client.Get(url)

			// Готовим результат для отправки в канал.
			result := Result{URL: url}

			if err != nil {
				result.Err = err
				results <- result
				return
			}

			resp.Body.Close()
			result.OK = (resp.StatusCode == http.StatusOK)
			results <- result
		}()
	}

	// Собираем результаты из канала в основном потоке.
	// Знаем точное количество результатов = len(urls).
	for i := 0; i < len(urls); i++ {
		result := <-results
		if result.Err != nil {
			fmt.Printf("%s - not ok (error: %v)\n", result.URL, result.Err)
		} else if result.OK {
			fmt.Printf("%s - ok\n", result.URL)
		} else {
			fmt.Printf("%s - not ok\n", result.URL)
		}
	}
}

// ============================================================================
// Этап 3: Без использования длины слайса
// ============================================================================

func step3Unknown() {
	fmt.Println("\n=== Этап 3: Неизвестное количество URL ===")

	// Симулируем источник URL'ов через канал.
	// В реальности это может быть очередь, стрим, сокет и т.д.
	urlChan := make(chan string)

	// Горутина, генерирующая URL'ы.
	go func() {
		var urls = []string{
			"http://ozon.ru",
			"https://ozon.ru",
			"http://google.com",
			"http://somesite.com",
			"http://non-existent.domain.tld",
			"https://ya.ru",
			"http://ya.ru",
			"http://eeeboy",
		}

		for _, url := range urls {
			urlChan <- url
		}
		// ВАЖНО: закрываем канал, сигнализируя об окончании данных.
		// Это идиоматичный способ в Go.
		close(urlChan)
	}()

	client := &http.Client{Timeout: 5 * time.Second}
	// Канал для результатов (без буфера, синхронная отправка).
	results := make(chan Result)

	// WaitGroup для ожидания завершения всех горутин-воркеров.
	var wg sync.WaitGroup

	// Запускаем воркеры для обработки URL'ов из канала.
	for url := range urlChan {
		wg.Add(1)
		url := url

		go func() {
			defer wg.Done()

			resp, err := client.Get(url)

			result := Result{URL: url}
			if err != nil {
				result.Err = err
			} else {
				resp.Body.Close()
				result.OK = (resp.StatusCode == http.StatusOK)
			}

			results <- result
		}()
	}

	// Горутина для закрытия канала результатов после завершения всех воркеров.
	go func() {
		wg.Wait()      // Ждем завершения всех воркеров
		close(results) // Закрываем канал результатов
	}()

	// Читаем результаты до закрытия канала.
	// range автоматически завершится, когда канал закроется.
	for result := range results {
		if result.Err != nil {
			fmt.Printf("%s - not ok (error: %v)\n", result.URL, result.Err)
		} else if result.OK {
			fmt.Printf("%s - ok\n", result.URL)
		} else {
			fmt.Printf("%s - not ok\n", result.URL)
		}
	}
}

// ============================================================================
// Этап 4: Отмена через context после 2 успешных ответов
// ============================================================================

func step4Context() {
	fmt.Println("\n=== Этап 4: Отмена через context ===")

	var urls = []string{
		"http://ozon.ru",
		"https://ozon.ru",
		"http://google.com",
		"http://somesite.com",
		"http://non-existent.domain.tld",
		"https://ya.ru",
		"http://ya.ru",
		"http://eeeboy",
	}

	client := &http.Client{Timeout: 5 * time.Second}
	// Создаем context с возможностью отмены.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Гарантируем освобождение ресурсов

	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	// Запускаем горутины для запросов.
	for _, url := range urls {
		wg.Add(1)
		url := url

		go func() {
			defer wg.Done()

			// Создаем запрос с контекстом.
			req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
			if err != nil {
				results <- Result{URL: url, Err: err}
				return
			}

			// Выполняем запрос с таймаутом и контекстом.
			resp, err := client.Do(req)

			result := Result{URL: url}

			// Проверяем, был ли запрос отменен через контекст.
			if ctx.Err() != nil {
				fmt.Printf("%s - canceled\n", url)
				return
			}

			if err != nil {
				result.Err = err
			} else {
				resp.Body.Close()
				result.OK = (resp.StatusCode == http.StatusOK)
			}

			results <- result
		}()
	}

	// Горутина для закрытия канала результатов.
	go func() {
		wg.Wait()
		close(results)
	}()

	// Счетчик успешных ответов.
	successCount := 0

	// Читаем результаты и отменяем контекст после 2 успешных.
	for result := range results {
		if result.Err != nil {
			fmt.Printf("%s - not ok (error: %v)\n", result.URL, result.Err)
		} else if result.OK {
			fmt.Printf("%s - ok\n", result.URL)
			successCount++

			// После 2 успешных ответов отменяем контекст.
			if successCount >= 2 {
				fmt.Println("Получено 2 успешных ответа, отменяем оставшиеся запросы")
				cancel()
			}
		} else {
			fmt.Printf("%s - not ok\n", result.URL)
		}
	}
}

// ============================================================================
// Этап 5: Рефакторинг и тестируемость
// ============================================================================

// HTTPClient интерфейс для абстракции HTTP клиента (для тестирования).
type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

// URLChecker проверяет доступность URL'ов.
type URLChecker struct {
	client  HTTPClient
	timeout time.Duration
}

// NewURLChecker создает новый checker с настройками.
func NewURLChecker(client HTTPClient, timeout time.Duration) *URLChecker {
	return &URLChecker{
		client:  client,
		timeout: timeout,
	}
}

// Check проверяет один URL с учетом контекста.
func (c *URLChecker) Check(ctx context.Context, url string) Result {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return Result{URL: url, Err: err}
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return Result{URL: url, Err: err}
	}
	defer resp.Body.Close()

	return Result{
		URL: url,
		OK:  resp.StatusCode == http.StatusOK,
	}
}

// CheckMany проверяет несколько URL'ов конкурентно.
func (c *URLChecker) CheckMany(ctx context.Context, urls []string) <-chan Result {
	results := make(chan Result, len(urls))
	var wg sync.WaitGroup

	for _, url := range urls {
		wg.Add(1)
		url := url

		go func() {
			defer wg.Done()
			results <- c.Check(ctx, url)
		}()
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	return results
}

func main() {
	// Демонстрация всех этапов
	step1Sequential()
	step2Concurrent()
	step3Unknown()
	step4Context()

	// Пример использования рефакторенного кода
	fmt.Println("\n=== Этап 5: Рефакторинг (пример) ===")
	checker := NewURLChecker(&http.Client{}, 5*time.Second)
	ctx := context.Background()

	urls := []string{"https://google.com", "https://ya.ru"}
	results := checker.CheckMany(ctx, urls)

	for result := range results {
		if result.OK {
			fmt.Printf("%s - ok\n", result.URL)
		} else {
			fmt.Printf("%s - not ok\n", result.URL)
		}
	}
}

// Объяснение решения:
//
// 1. Этап 1 - Последовательные запросы:
//    - Простой цикл по URL'ам
//    - Синхронное выполнение http.Get
//    - Таймаут для предотвращения зависания
//    - Важно закрывать resp.Body
//
// 2. Этап 2 - Конкурентность с каналами:
//    - Запуск горутины для каждого URL
//    - Буферизированный канал для результатов
//    - Важно захватывать url в локальную переменную для горутины
//    - Основной поток читает из канала и печатает
//
// 3. Этап 3 - Неизвестное количество:
//    - Источник URL'ов - входной канал
//    - Закрытие канала сигнализирует об окончании
//    - WaitGroup для ожидания завершения воркеров
//    - range по каналу автоматически завершается при закрытии
//    - Паттерн: генератор -> воркеры -> сборщик результатов
//
// 4. Этап 4 - Context для отмены:
//    - context.WithCancel для управления жизненным циклом
//    - http.NewRequestWithContext связывает запрос с контекстом
//    - cancel() отменяет все активные запросы
//    - Проверка ctx.Err() для определения отмены
//    - Счетчик успешных ответов для логики отмены
//
// 5. Этап 5 - Тестируемость:
//    - HTTPClient интерфейс для моков
//    - URLChecker инкапсулирует логику проверки
//    - Dependency injection через конструктор
//    - Методы принимают context для управления
//    - Легко писать unit-тесты с mock клиентом
//
// Как тестировать:
//
// type mockHTTPClient struct {
//     response *http.Response
//     err      error
// }
//
// func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
//     return m.response, m.err
// }
//
// func TestURLChecker_Check(t *testing.T) {
//     mock := &mockHTTPClient{
//         response: &http.Response{StatusCode: 200},
//     }
//     checker := NewURLChecker(mock, 5*time.Second)
//     result := checker.Check(context.Background(), "http://test.com")
//     if !result.OK {
//         t.Error("expected OK")
//     }
// }
//
// Trade-offs:
// - Последовательно: медленно, но просто
// - Конкурентно: быстро, но сложнее отладка
// - С context: гибкость отмены, но больше кода
// - С интерфейсами: тестируемость, но больше абстракций
