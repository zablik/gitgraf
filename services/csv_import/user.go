package csvimport

import (
	"context"
	"encoding/json"
	"gitgraf/config"
	"gitgraf/model"
	"gitgraf/repository"
	"gitgraf/repository/mongodb"
	"gitgraf/services/utils"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	userRepoObj repository.UserRepository
)

func GetPatch() map[string]interface{} {
	if len(os.Args) < 2 {
		log.Fatal("Error: Missing command-line argument")
	}

	base := filepath.Base(os.Args[1])
	name := strings.TrimSuffix(base, ".csv")
	filename := filepath.Join(filepath.Dir(os.Args[1]), name+"-patch.json")

	_, err := os.Stat(filename)
	if err != nil {
		log.Fatal("Error:", err)
	}

	// Read the JSON file
	data, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal("Failed to read JSON file:", err)
	}

	// Define a map to store the JSON data
	var jsonData map[string]interface{}

	// Unmarshal the JSON data into the map
	err = json.Unmarshal(data, &jsonData)
	if err != nil {
		log.Fatal("Failed to unmarshal JSON:", err)
	}

	return jsonData
}

func SaveUsers(ctx context.Context, client *mongo.Client, users map[string]model.User, wg *sync.WaitGroup) {
	_, err := userRepo(client).CreateMany(ctx, ConvertMapToUserSlice(users))
	if err != nil {
		log.Printf("Error inserting documents: %v", err)
	}
	wg.Done()
}

func ProcessUsers(users *map[string]model.User, lines <-chan []string, usersLock *sync.Mutex, numWorkers int) {
	var wg sync.WaitGroup
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for line := range lines {
				email := line[1]
				name := line[2]
				datetime, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", line[3])
				if err != nil {
					log.Fatal("Failed to parse datetime:", line[3], " ===> ", err)
					continue
				}

				usersLock.Lock()
				user, ok := (*users)[email]
				if !ok {
					user = model.User{
						Email:         email,
						Name:          name,
						AltNames:      []string{},
						FirstActiveAt: datetime,
						LastActiveAt:  datetime,
					}
				} else {
					if user.Name != name && !utils.Contains(user.AltNames, name) {
						user.AltNames = append(user.AltNames, name)
					}

					if user.LastActiveAt.Before(datetime) {
						user.LastActiveAt = datetime
						if user.Name != name {
							user.Name = name
						}
					}

					if user.FirstActiveAt.After(datetime) {
						user.FirstActiveAt = datetime
					}
				}
				(*users)[email] = user
				usersLock.Unlock()
			}
		}()
	}
	wg.Wait()
}

func ApplyPatch(users *map[string]model.User, patch map[string]interface{}) {
	for email, altEmails := range patch["email"].(map[string]interface{}) {
		mainUser, ok := (*users)[email]
		if !ok {
			log.Println("User not found:", email)
			continue
		}

		for _, altEmail := range altEmails.([]interface{}) {
			altUser, ok := (*users)[altEmail.(string)]
			if !ok {
				log.Println("Alt User not found:", altEmail)
				continue
			}

			mainUser.AltEmails = append(mainUser.AltEmails, altEmail.(string))

			if !utils.Contains(mainUser.AltNames, altUser.Name) {
				mainUser.AltNames = append(mainUser.AltNames, altUser.Name)
			}

			if mainUser.LastActiveAt.Before(altUser.LastActiveAt) {
				mainUser.LastActiveAt = altUser.LastActiveAt
			}

			if mainUser.FirstActiveAt.After(altUser.FirstActiveAt) {
				mainUser.FirstActiveAt = altUser.FirstActiveAt
			}

			// Delete an item from the map
			delete((*users), altEmail.(string))
		}

		(*users)[email] = mainUser
	}
}

func userRepo(client *mongo.Client) repository.UserRepository {
	if userRepoObj == nil {
		userRepoObj = mongodb.NewUserRepository(client, config.Load().DB.Name, "users")
	}
	return userRepoObj
}

func AltEmailMap(ctx context.Context, client *mongo.Client) map[string]*model.User {
	emailMap := make(map[string]*model.User)

	allUsers, err := userRepo(client).GetAll(ctx)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
	}

	for _, user := range allUsers {
		for _, altEmail := range user.AltEmails {
			emailMap[altEmail] = user
		}
		emailMap[user.Email] = user
	}

	return emailMap
}

func AltNameMap(ctx context.Context, client *mongo.Client) map[string]*model.User {
	nameMap := make(map[string]*model.User)

	allUsers, err := userRepo(client).GetAll(ctx)
	if err != nil {
		log.Printf("Error getting all users: %v", err)
	}

	for _, user := range allUsers {
		for _, altName := range user.AltNames {
			nameMap[altName] = user
		}
	}

	return nameMap
}

func ConvertMapToUserSlice(m map[string]model.User) []*model.User {
	var s []*model.User
	for _, v := range m {
		s = append(s, &v)
	}
	return s
}
