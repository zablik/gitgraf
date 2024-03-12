package parsers

import (
	"gitgraf/model"
)

func parseCSVWorker(jobs <-chan []string, results chan<- (model.User, model.Commit)) {
	for record := range jobs {


		// Предположим, что первый столбец - имя, второй - возраст
		user := model.User{Name: record[0], Age: parseInt(record[1])}
		results <- user
	}
}

func lineToCommit(line string) (model.User, model.Commit) {

}
