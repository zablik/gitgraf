package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"
)

type User struct {
	Name string
	Age  int
}

func worker(jobs <-chan []string, results chan<- User, wg *sync.WaitGroup) {
	for record := range jobs {

		userAge, err := parseInt(record[1])
		if err != nil {
			log.Printf("Не удалось конвертировать возраст '%s' в число: %v", record[1], err)
			// Обработайте ошибку соответствующим образом
			continue // Например, пропустите текущую итерацию цикла
		}

		user := User{
			Name: record[0],
			Age:  userAge,
		}
		results <- user
	}
	wg.Done()
}

func main() {
	// Setup MongoDB connection...

	jobs := make(chan []string, 100) // Буферизованный канал для строк CSV
	results := make(chan User, 100)  // Буферизованный канал для обработанных данных

	// Инициализация пула воркеров
	var wg sync.WaitGroup
	for w := 1; w <= 10; w++ { // Количество воркеров
		wg.Add(1)
		go worker(jobs, results, &wg)
	}

	// Чтение и отправка данных в канал jobs
	go func() {
		csvFile, err := os.Open("data.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer csvFile.Close()

		reader := csv.NewReader(bufio.NewReader(csvFile))
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal(err)
			}
			jobs <- record
		}
		close(jobs)
	}()

	// Обработка результатов асинхронно
	go func() {
		for user := range results {
			// Здесь можно сохранить user в MongoDB
			fmt.Println("Processed user:", user)
		}
	}()

	wg.Wait()
	close(results)
	// Закрытие соединения с MongoDB и другие операции по завершению...
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseInt и другие необходимые функции...
