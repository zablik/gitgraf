package cmd

import (
	"encoding/csv"
	"fmt"
	"gitgraf/model"
	"io"
	"log"
	"os"
)

func userCsvToDb() {

	csvFile, err := os.Open("path/to/your/file.csv")
	if err != nil {
		log.Fatal("Failed to open CSV file:", err)
	}
	defer csvFile.Close()

	var csvUsers []model.User

	reader := csv.NewReader(csvFile)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Failed to read CSV line:", err)
		}
		csvUsers = append(csvUsers, convertCsvUserLineToUser(record))
	}

}

func findSameEmail(users []model.User) {
	for i, user1 := range users {
		for j, user2 := range users {
			if i >= j { // Избегаем повторной проверки и самопроверки
				continue
			}
			if user1.Email == user2.Email {
				user2.Main


		}
	}
}

func convertCsvUserLineToUser(lineChunks []string) model.User {
	return model.User{
		Name:  lineChunks[0],
		Email: lineChunks[1],
		Main:  nil,
	}
}
