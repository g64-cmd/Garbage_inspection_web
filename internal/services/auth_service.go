package services

import (
	"context"
	"errors"
	"patrol-cloud/internal/db"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

// AuthService 提供了用户认证相关的服务
type AuthService struct {
	repo         db.Repository
	jwtSecretKey []byte
}

// NewAuthService 创建一个新的 AuthService
func NewAuthService(repo db.Repository, jwtSecret []byte) *AuthService {
	return &AuthService{repo: repo, jwtSecretKey: jwtSecret}
}

// Login 验证用户凭据并返回一个 JWT
func (s *AuthService) Login(ctx context.Context, username, password string) (string, error) {
	// 1. 从数据库获取用户
	user, err := s.repo.GetUserByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", ErrUserNotFound
	}

	// 2. 比较密码哈希
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		// 如果密码不匹配，bcrypt 会返回一个错误
		return "", ErrInvalidCredentials
	}

	// 3. 生成 JWT
	token, err := s.generateJWT(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// generateJWT 为指定用户生成一个新的 JWT
func (s *AuthService) generateJWT(userID, username string) (string, error) {
	// 创建 claims
	claims := jwt.MapClaims{
		"sub":      userID,
		"username": username,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // 24小时后过期
	}

	// 创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名并获取完整的 token 字符串
	tokenString, err := token.SignedString(s.jwtSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 解析并验证一个 JWT 字符串
func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// 确保签名方法是我们期望的
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecretKey, nil
	})
}
