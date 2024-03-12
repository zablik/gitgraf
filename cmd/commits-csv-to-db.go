package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"sync"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	workersCount = 10  // Количество воркеров для парсинга CSV
	batchSize    = 100 // Размер пакета для записи в MongoDB
)

// type User struct {
// 	Name string
// 	Age  int
// }

func parseCSVWorker(jobs <-chan []string, results chan<- User) {
	for record := range jobs {
		// Предположим, что первый столбец - имя, второй - возраст
		user := User{Name: record[0], Age: parseInt(record[1])}
		results <- user
	}
}

func batchInsertUsers(ctx context.Context, client *mongo.Client, users []User) {
	// Преобразование пользователей в []interface{} для InsertMany
	var documents []interface{}
	for _, user := range users {
		documents = append(documents, user)
	}

	collection := client.Database("yourDatabase").Collection("users")
	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Printf("Error inserting documents: %v", err)
	}
}

func main() {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(context.Background())

	csvFile, err := os.Open("path/to/your/file.csv")
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}
	defer csvFile.Close()

	jobs := make(chan []string)
	results := make(chan User, batchSize)

	// Запуск воркеров для парсинга CSV
	for i := 0; i < workersCount; i++ {
		go parseCSVWorker(jobs, results)
	}

	// Чтение CSV и отправка строк в канал jobs
	go func() {
		reader := csv.NewReader(csvFile)
		for {
			record, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatal("Failed to read CSV line:", err)
			}
			jobs <- record
		}
		close(jobs)
	}()

	// Пакетная обработка и запись результатов
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		var batch []User
		for user := range results {
			batch = append(batch, user)
			if len(batch) >= batchSize {
				batchInsertUsers(context.Background(), client, batch)
				batch = nil // Сброс пакета после вставки
			}
		}
		// Обработка оставшихся пользователей
		if len(batch) > 0 {
			batchInsertUsers(context.Background(), client, batch)
		}
	}()

	wg.Wait()
	fmt.Println("Data processing completed.")
}

func parseInt(s string) (int, error) {
	return strconv.Atoi(s)
}

// parseInt и другие необходимые функции...
