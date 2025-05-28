package permission

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Permission struct {
	ID           primitive.ObjectID `bson:"id"`            // MongoDB ObjectID
	UuidId       string             `bson:"uuid_id"`       // Дополнительный UUID идентификатор
	ResourceType string             `bson:"resource_type"` // Тип ресурса: "document", "folder", "file"
	ResourceID   primitive.ObjectID `bson:"resource_id"`   // ID ресурса в MongoDB
	ResourceUuid string             `bson:"resource_uuid"` // UUID ресурса (альтернативный идентификатор)
	Rules        PermissionRules    `bson:"rules"`         // Правила доступа
	CreatedAt    time.Time          `bson:"created_at"`    // Время создания
	UpdatedAt    time.Time          `bson:"updated_at"`    // Время последнего обновления
	CreatedBy    string             `bson:"created_by"`    // Кто создал (user_id/uuid)
	Version      string             `bson:"version"`       // Версия для оптимистичной блокировки
}

type PermissionRules struct {
	AccessAllowed []string `bson:"access_allowed"` // Список ID пользователей с разрешенным доступом
	CommentOnly   []string `bson:"comment_only"`   // Список ID пользователей с доступом только для комментирования
	ReadOnly      []string `bson:"read_only"`      // Список ID пользователей с доступом только для чтения
	AccessLevel   string   `bson:"access_level"`   // Общий уровень доступа: "public", "private", "restricted"
}

func NewPermission(employeeId, generatedId, generatedResourceUuid, resourceType string, resourceId primitive.ObjectID) Permission {
	return Permission{
		ID:           primitive.NewObjectID(),
		UuidId:       generatedId,
		ResourceType: resourceType,
		ResourceID:   resourceId,
		ResourceUuid: generatedId,
		Rules:        NewPermissionRules(employeeId),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		CreatedBy:    employeeId,
		Version:      "1.0.0",
	}
}

func NewPermissionRules(employeeId string) PermissionRules {
	return PermissionRules{
		[]string{employeeId},
		[]string{},
		[]string{},
		"",
	}
}
