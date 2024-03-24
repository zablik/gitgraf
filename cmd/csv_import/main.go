package main

import (
	"context"
	"encoding/csv"
	"fmt"
	"gitgraf/config"
	"gitgraf/db"
	"gitgraf/model"
	csvimport "gitgraf/services/csv_import"
	"gitgraf/services/utils"

	"io"
	"log"
	"os"
	"sync"
)

var (
	users       = make(map[string]model.User)
	usersLock   sync.Mutex
	commitsLock sync.Mutex
)

func main() {
	ctx := context.TODO()
	cfg := config.Load()
	client, err := db.SetupMongo(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	file := getFile()

	defer func() {
		client.Disconnect(context.Background())
		file.Close()
	}()

	lines := make(chan []string, 100)

	utils.MeasureTime("Process Users", func() {
		var wg sync.WaitGroup
		wg.Add(1)
		go csvimport.ProcessUsers(&users, lines, &usersLock, 3)
		readFile(lines, file)
		csvimport.ApplyPatch(&users, csvimport.GetPatch())
		csvimport.SaveUsers(ctx, client, users, &wg)

		wg.Wait()
	})

	for _, value := range users {
		fmt.Println(value.Email, value.Name, value.AltEmails, value.AltNames)
	}

	lines = make(chan []string, 100)
	commits := make(chan model.Commit, 100)

	utils.MeasureTime("Process Commits", func() {
		var wg sync.WaitGroup
		go csvimport.ProcessCommits(
			lines,
			commits,
			csvimport.AltEmailMap(ctx, client),
			csvimport.AltNameMap(ctx, client),
			&commitsLock,
			3,
		)
		go csvimport.SaveCommits(ctx, client, commits, &commitsLock, &wg)

		readFile(lines, getFile())

		wg.Wait()
	})
}

func funkYou() {
	fmt.Println("funkYou")
}

func readFile(lines chan<- []string, file *os.File) {
	reader := csv.NewReader(file)
	isFirstLine := true
	defer close(lines)

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Failed to read CSV line:", err)
		}

		if isFirstLine {
			isFirstLine = false
			continue // Skip the first line
		}

		lines <- record
	}
}

func getFile() *os.File {
	if len(os.Args) < 2 {
		log.Fatal("Error: Missing command-line argument")
	}

	filename := os.Args[1]
	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal("Error:", err)
	}
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}

	return file
}
