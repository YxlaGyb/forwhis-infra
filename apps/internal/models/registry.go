package models

import "forestwhisperGo/apps/internal/models/rbac"

func All() []any {
	return []any{
		&User{},
		&Userinfo{},
		&Category{},
		&Tag{},
		&Content{},
		&ContentTag{},
		&Comment{},
		&Like{},
		&Option{},
		&Review{},
		&AIConversation{},
		&rbac.UserRole{},
	}
}
