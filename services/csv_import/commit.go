package csvimport

import (
	"context"
	"encoding/json"
	"fmt"
	"gitgraf/config"
	"gitgraf/model"
	"gitgraf/repository"
	"gitgraf/repository/mongodb"
	"log"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	commitRepoObj      repository.CommitRepository
)

func ProcessCommits(
	lines <-chan []string,
	commits chan<- model.Commit,
	altEmailMap map[string]*model.User,
	altNameMap map[string]*model.User,
	commitsLock *sync.Mutex,
	workers int,
) {
	var wg sync.WaitGroup
	defer close(commits)
	
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for line := range lines {
				processCommitLine(line, commits, altEmailMap, altNameMap)
			}
		}()
	}
	wg.Wait()
}

func SaveCommits(ctx context.Context, client *mongo.Client, commits <-chan model.Commit, commitsLock *sync.Mutex, wg *sync.WaitGroup) {
	batchSize := 5000
	var batch []*model.Commit
	wg.Add(1)

	for commit := range commits {
		batch = append(batch, &commit)

		if len(batch) == batchSize {
			saveBatch(batch, ctx, client, commitsLock)
			batch = nil
		}
	}

	// Save any remaining commits in the batch
	if len(batch) > 0 {
		saveBatch(batch, ctx, client, commitsLock)
	}

	wg.Done()
}

func saveBatch(commits []*model.Commit, ctx context.Context, client *mongo.Client, commitsLock *sync.Mutex) {
	commitsLock.Lock()
	_, err := commitRepo(client).CreateMany(ctx, commits)
	if err != nil {
		log.Printf("Error inserting documents: %v", err)
	}
	commitsLock.Unlock()
}

func processCommitLine(
	line []string,
	commits chan<- model.Commit,
	altEmailMap map[string]*model.User,
	altNameMap map[string]*model.User,
) {
	// ID        primitive.ObjectID   `bson:"_id,omitempty"`
	// UserId    primitive.ObjectID   `bson:"user_id"`
	// CreatedAt time.Time            `bson:"created_at"`
	// Approvers []primitive.ObjectID `bson:"approvers"`
	// Stats     Stats                `bson:"stats"`

	// LinesAdded    int `bson:"lines_added"`
	// LinesDeleted  int `bson:"lines_deleted"`
	// FilesAdded    int `bson:"files_added"`
	// FilesDeleted  int `bson:"files_deleted"`
	// FilesModified int `bson:"files_modified"`

	// 0 - hash
	// 1 - email
	// 2 - name
	// 3 - date
	// 4 - files_changed
	// 5 - files_created
	// 6 - files_deleted
	// 7 - lines_added
	// 8 - lines_deleted
	// 9 - approvers

	userId := altEmailMap[line[1]].ID

	approvers := []string{}

	if err := json.Unmarshal([]byte(line[9]), &approvers); err != nil {
		fmt.Println("Error unmarshalling contacts:", err)
	}

	approverIds := []primitive.ObjectID{}
	for _, approverName := range approvers {
		if altNameMap[approverName] != nil {
			approverId := altNameMap[approverName].ID
			approverIds = append(approverIds, approverId)
		}
	}

	createdAt, err := time.Parse("Mon Jan 2 15:04:05 2006 -0700", line[3])
	if err != nil {
		fmt.Println("Failed to parse datetime:", line[3], " ===> ", err)
	}

	linesAdded, err := strconv.Atoi(line[7])
	if err != nil {
		fmt.Println("Failed to parse lines added:", line[7], " ===> ", err)
	}

	linesDeleted, err := strconv.Atoi(line[8])
	if err != nil {
		fmt.Println("Failed to parse lines deleted:", line[8], " ===> ", err)
	}

	filesAdded, err := strconv.Atoi(line[5])
	if err != nil {
		fmt.Println("Failed to parse files added:", line[5], " ===> ", err)
	}

	filesDeleted, err := strconv.Atoi(line[6])
	if err != nil {
		fmt.Println("Failed to parse files deleted:", line[6], " ===> ", err)
	}

	filesModified, err := strconv.Atoi(line[4])
	if err != nil {
		fmt.Println("Failed to parse files modified:", line[4], " ===> ", err)
	}

	commits <- model.Commit{
		UserId:      userId,
		Hash:        line[0],
		CreatedAt:   createdAt,
		ApproverIds: approverIds,
		Stats: model.Stats{
			LinesAdded:    linesAdded,
			LinesDeleted:  linesDeleted,
			FilesAdded:    filesAdded,
			FilesDeleted:  filesDeleted,
			FilesModified: filesModified,
		},
	}
}

func ConvertMapToCommitSlice(m map[string]model.Commit) []*model.Commit {
	var s []*model.Commit
	for _, v := range m {
		s = append(s, &v)
	}
	return s
}

func commitRepo(client *mongo.Client) repository.CommitRepository {
	if commitRepoObj == nil {
		commitRepoObj = mongodb.NewCommitRepository(client, config.Load().DB.Name, "commits")
	}
	return commitRepoObj
}
