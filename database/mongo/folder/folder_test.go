package folder_test

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/database/mongo/folder"
	"labyrinth/notebook/models/directory"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client        *mongo.Client
	testDB        *mongo.Database
	testDirectory *directory.Directory
)

func setup() error {
	var err error
	client, err = m.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	testDB = client.Database("folder_test")

	testDirectory = &directory.Directory{
		ID:        primitive.NewObjectID(),
		UuidID:    uuid.New().String(),
		ParentId:  uuid.New().String(),
		IsPrimary: true,
		Version:   "1.0",
		Metadata: directory.Metadata{
			CompanyID:   uuid.New().String(),
			DivisionID:  uuid.New().String(),
			Title:       "TITLE",
			Description: "DESCRIPTION",
			Tags:        []string{"INFO", "SCIENCE"},
			Created: directory.Timestamp{
				Date:   time.Now(),
				Time:   time.Now(),
				Author: uuid.New().String(),
			},
			LastUpdate: directory.Timestamp{
				Date:   time.Now(),
				Time:   time.Now(),
				Author: uuid.New().String(),
			},
			Links: directory.DocumentLinks{
				Read:    "https://api.example.com/v1/documents/123/read",
				Comment: "https://api.example.com/v1/documents/123/comments",
				Write:   "https://api.example.com/v1/documents/123/write",
				ActiveLinks: []string{
					"read",
					"comment",
				},
			},
		},
		Folders: []directory.Folder{
			{
				FolderID:    primitive.NewObjectID(),
				FolderUUID:  uuid.New().String(),
				Title:       "FOLDER TITLE 1",
				Description: "FOLDER DESCRIPTION 1",
			},
			{
				FolderID:    primitive.NewObjectID(),
				FolderUUID:  uuid.New().String(),
				Title:       "FOLDER TITLE 2",
				Description: "FOLDER DESCRIPTION 2",
			},
		},
		Files: []directory.File{
			{
				FileID:      primitive.NewObjectID(),
				FileUUID:    uuid.New().String(),
				Title:       "FILE TITLE 1",
				Description: "FILE DESCRIPTION 1",
			},
		},
	}
	return nil
}

func teardown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if testDB != nil {
		testDB.Drop(ctx)
	}

	if client != nil {
		client.Disconnect(ctx)
	}
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		fmt.Printf("Test setup failed: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	teardown()

	os.Exit(code)
}

func TestFolderCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := folder.NewFolderMongo(testDB, "folder_test")

	session, err := client.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		t.Fatalf("Failed to start transaction: %v\n", err)
	}

	t.Run("CreateFolder", func(t *testing.T) {
		err = repo.CreateFolder(ctx, &session, testDirectory)
		if err != nil {
			t.Fatalf("CreateFolder failed %v\n", err)
		}
	})

	t.Run("ExistsFolder", func(t *testing.T) {
		exists, err := repo.ExistsFolder(ctx, &session, testDirectory.UuidID)
		if err != nil {
			t.Fatalf("ExistsFolder failed: %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got false\n")
		}
	})

	t.Run("UpdateFolder", func(t *testing.T) {
		updatedDirectory := *testDirectory
		updatedDirectory.Folders = append(updatedDirectory.Folders,
			directory.Folder{
				FolderID:    primitive.NewObjectID(),
				FolderUUID:  uuid.New().String(),
				Title:       "FOLDER TITLE 3",
				Description: "FOLDER DESCRIPTION 3",
			})
		updatedDirectory.Metadata.LastUpdate = directory.Timestamp{
			Date:   time.Now(),
			Time:   time.Now(),
			Author: uuid.New().String(),
		}

		err = repo.UpdateFolder(ctx, &session, updatedDirectory.UuidID, &updatedDirectory)
		if err != nil {
			t.Fatalf("UpdateFolder faield: %v\n", err)
		}
	})

	t.Run("GetFoldersByParentId", func(t *testing.T) {
		fetchedFolders, err := repo.GetFoldersByParentId(ctx, &session, testDirectory.ParentId)
		if err != nil {
			t.Fatalf("GetFoldersByParentId failed: %v\n", err)
		}

		if len(fetchedFolders) != 1 {
			t.Errorf("Expected len 1, got %d\n", len(fetchedFolders))
		}

		for _, value := range fetchedFolders {
			if len(value.Folders) != 3 {
				t.Errorf("Expected len 3, got %d\n", len(value.Folders))
			}
		}
	})

	t.Run("DeleteFolder", func(t *testing.T) {
		err = repo.DeleteFolder(ctx, &session, testDirectory.UuidID)
		if err != nil {
			t.Fatalf("DeleteFolder failed: %v\n", err)
		}
	})

	if !t.Failed() {
		if err := session.CommitTransaction(ctx); err != nil {
			t.Errorf("Failed to commit transaction: %v", err)
		}
	} else {
		if err := session.AbortTransaction(ctx); err != nil {
			t.Errorf("Failed to abort transaction: %v", err)
		}
	}
}
