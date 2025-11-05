//go:build dev

package main

import (
	"context"
	"log"
	"patrol-cloud/internal/db"
	"patrol-cloud/internal/models"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// createTempUser 用于在启动时创建一个测试用户 (仅在开发构建中包含)
func createTempUser(repo db.Repository) {
	ctx := context.Background()
	username := "admin"

	// 检查用户是否已存在
	existingUser, err := repo.GetUserByUsername(ctx, username)
	if err != nil {
		log.Printf("WARN: Could not check for temp user: %v", err)
		return
	}
	if existingUser != nil {
		log.Println("INFO: Temp user 'admin' already exists.")
		return
	}

	// 创建用户
	password := "password123"
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("WARN: Could not hash password for temp user: %v", err)
		return
	}

	user := &models.User{
		ID:             uuid.NewString(),
		Username:       username,
		HashedPassword: string(hashedPassword),
	}

	if err := repo.CreateUser(ctx, user); err != nil {
		log.Printf("WARN: Could not create temp user: %v", err)
		return
	}

	log.Printf("INFO: Created temp user '%s' with password '%s'", username, password)
}
