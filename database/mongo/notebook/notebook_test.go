package notebook_test

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/database/mongo/notebook"
	"labyrinth/notebook/models/journal"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client       *mongo.Client
	testDB       *mongo.Database
	testNotebook *journal.Notebook
)

func setup() error {
	var err error
	client, err = m.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	testDB = client.Database("notebook_test")

	testNotebook = &journal.Notebook{
		ID:      primitive.NewObjectID(),
		UuidID:  uuid.New().String(),
		Version: "1.0",
		Metadata: journal.Metadata{
			CompanyID:   uuid.New().String(),
			DivisionID:  uuid.New().String(),
			Title:       "TITLE",
			Description: "DESCRIPTION",
			Tags:        []string{"INFO", "SCIENCE"},
			Created: journal.DateTimeAuthor{
				Date:   time.Now(),
				Time:   time.Now(),
				Author: uuid.New().String(),
			},
			LastUpdate: journal.DateTimeAuthor{
				Date:   time.Now(),
				Time:   time.Now(),
				Author: uuid.New().String(),
			},
			Links: journal.Links{
				Read:    "https://api.example.com/v1/documents/123/read",
				Comment: "https://api.example.com/v1/documents/123/comments",
				Write:   "https://api.example.com/v1/documents/123/write",
				ActiveLinks: []string{
					"read",
					"comment",
				},
			},
		},
		Blocks: []journal.Block{
			{
				Id:   1,
				Type: "title",
				Body: map[string]any{
					"text": "Заголовок документа",
				},
				Comment: []journal.Comment{
					{
						EmployeeId: "emp-001",
						CreatedAt:  time.Now().Add(-48 * time.Hour),
						UpdatedAt:  time.Now().Add(-24 * time.Hour),
						Comment:    "Нужно добавить раздел по маркетингу",
						SubComment: []journal.Comment{
							{
								EmployeeId: "emp-002",
								CreatedAt:  time.Now().Add(-12 * time.Hour),
								UpdatedAt:  time.Now(),
								Comment:    "Добавил раздел, проверьте",
								SubComment: []journal.Comment{},
							},
						},
					},
					{
						EmployeeId: "emp-003",
						CreatedAt:  time.Now().Add(-1 * time.Hour),
						UpdatedAt:  time.Now(),
						Comment:    "Поправил заголовок",
						SubComment: []journal.Comment{},
					},
				},
			},
			{
				Id:   2,
				Type: "text",
				Body: map[string]any{
					"content": "Основные показатели за год...",
				},
				Comment: []journal.Comment{
					{
						EmployeeId: "emp-004",
						CreatedAt:  time.Now().Add(-3 * time.Hour),
						UpdatedAt:  time.Now().Add(-1 * time.Hour),
						Comment:    "Нужно уточнить цифры",
						SubComment: []journal.Comment{},
					},
				},
			},
			{
				Id:   3,
				Type: "table",
				Body: map[string]any{
					"columns": []string{"Показатель", "Значение"},
					"data": [][]string{
						{"Выручка", "1,200,000"},
						{"Прибыль", "350,000"},
					},
				},
				Comment: []journal.Comment{},
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

func TestNotebookCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := notebook.NewNotebookMongo(testDB, "notebook_test")

	session, err := client.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		t.Fatalf("Failed to start transaction: %v\n", err)
	}

	t.Run("CreateNotebook", func(t *testing.T) {
		err := repo.CreateNotebook(ctx, &session, testNotebook)
		if err != nil {
			t.Fatalf("CreateNotebook failed: %v\n", err)
		}
	})

	t.Run("ExistsNotebook", func(t *testing.T) {
		exists, err := repo.ExistsNotebook(ctx, &session, testNotebook.UuidID)
		if err != nil {
			t.Fatalf("ExistsNotebook failed: %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got false\n")
		}
	})

	t.Run("UpdateNotebook", func(t *testing.T) {
		updatedNotebook := *testNotebook
		updatedNotebook.Metadata.Title = "UPDATED TITLE"
		updatedNotebook.Metadata.Description = "UPDATED  DESCRIPTION"

		err := repo.UpdateNotebook(ctx, &session, testNotebook.UuidID, &updatedNotebook)
		if err != nil {
			t.Fatalf("UpdateNotebook failed: %v\n", err)
		}
	})

	t.Run("GetNotebookById", func(t *testing.T) {
		fetchedNotebook, err := repo.GetNotebookById(ctx, &session, testNotebook.UuidID)
		if err != nil {
			t.Fatalf("GetNotebookById failed: %v\n", err)
		}

		if fetchedNotebook.Metadata.Title != "UPDATED TITLE" {
			t.Errorf("Expected UPDATED TITLE, got  %s\n", fetchedNotebook.Metadata.Title)
		}
	})

	t.Run("DeleteNotebook", func(t *testing.T) {
		err := repo.DeleteNotebook(ctx, &session, testNotebook.UuidID)
		if err != nil {
			t.Fatalf("DeleteNotebook failed %v\n", err)
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
