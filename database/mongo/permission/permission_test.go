package permission_test

import (
	"context"
	"fmt"
	m "labyrinth/database/mongo"
	"labyrinth/database/mongo/permission"
	p "labyrinth/notebook/models/permission"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	client         *mongo.Client
	testDB         *mongo.Database
	testPermission *p.Permission
)

func setup() error {
	var err error
	client, err = m.NewConnection()
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	testDB = client.Database("permission_test")

	testPermission = &p.Permission{
		ID:           primitive.NewObjectID(),
		UuidId:       uuid.New().String(),
		ResourceType: "folder",
		ResourceID:   primitive.NewObjectID(),
		ResourceUuid: uuid.New().String(),
		Rules: p.PermissionRules{
			AccessAllowed: []string{uuid.New().String()},
			CommentOnly:   []string{uuid.New().String()},
			ReadOnly:      []string{uuid.New().String()},
			AccessLevel:   "private",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Version:   "1.0",
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

func TestPermissionCRUD(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	repo := permission.NewPermissionMongo(testDB, "permission_test")

	session, err := client.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	defer session.EndSession(ctx)

	err = session.StartTransaction()
	if err != nil {
		t.Fatalf("Failed to start transaction: %v\n", err)
	}

	t.Run("CreatePermission", func(t *testing.T) {
		err = repo.CreatePermission(ctx, &session, testPermission)
		if err != nil {
			t.Fatalf("CreatePermission failed: %v\n", err)
		}
	})

	t.Run("ExistsPermission", func(t *testing.T) {
		exists, err := repo.ExistsPermission(ctx, &session, testPermission.UuidId)
		if err != nil {
			t.Fatalf("ExistsPermission failed: %v\n", err)
		}

		if !exists {
			t.Errorf("Expected true, got false\n")
		}
	})

	t.Run("UpdatePermission", func(t *testing.T) {
		updatedPermission := *testPermission
		updatedPermission.ResourceType = "file"
		updatedPermission.UpdatedAt = time.Now()

		err = repo.UpdatePermission(ctx, &session, updatedPermission.UuidId, &updatedPermission)
		if err != nil {
			t.Fatalf("UpdatePermission failed: %v\n", err)
		}
	})

	t.Run("GetPermissionByUuidId", func(t *testing.T) {
		fetchedPermission, err := repo.GetPermissionByUuidId(ctx, &session, testPermission.UuidId)
		if err != nil {
			t.Fatalf("GetPermissionByUuidId failed: %v\n", err)
		}

		if fetchedPermission.ResourceType != "file" {
			t.Errorf("Expected ( file ), got ( %s )'\n", fetchedPermission.ResourceType)
		}

		fmt.Printf("UUID-ID => get: %s, want: %s\n", fetchedPermission.UuidId, testPermission.UuidId)
	})

	t.Run("DeletePermission", func(t *testing.T) {
		err := repo.DeletePermission(ctx, &session, testPermission.UuidId)
		if err != nil {
			t.Fatalf("DeletePermission faield: %v\n", err)
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
