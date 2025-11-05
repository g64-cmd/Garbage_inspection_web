package services

import (
	"context"
	"patrol-cloud/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

// MockRepository is a mock type for the db.Repository interface
type MockRepository struct {
	mock.Mock
}

// GetUserByUsername is a mock implementation of the GetUserByUsername method
func (m *MockRepository) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

// LogDecision is a mock implementation (needed to satisfy the interface)
func (m *MockRepository) LogDecision(ctx context.Context, result *models.DecisionResult, imageURL string, metadata models.DecisionRequestMetadata) error {
	args := m.Called(ctx, result, imageURL, metadata)
	return args.Error(0)
}

// CreateUser is a mock implementation (needed to satisfy the interface)
func (m *MockRepository) CreateUser(ctx context.Context, user *models.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func TestAuthService_Login(t *testing.T) {
	// --- Setup ---
	jwtSecret := []byte("test-secret")
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	mockUser := &models.User{
		ID:             "user-123",
		Username:       "testuser",
		HashedPassword: string(hashedPassword),
	}

	// --- Test Cases ---
	tests := []struct {
		name          string
		username      string
		password      string
		setupMock     func(repo *MockRepository)
		expectError   bool
		expectedError error
	}{
		{
			name:     "Successful Login",
			username: "testuser",
			password: "password123",
			setupMock: func(repo *MockRepository) {
				repo.On("GetUserByUsername", mock.Anything, "testuser").Return(mockUser, nil)
			},
			expectError: false,
		},
		{
			name:     "User Not Found",
			username: "nonexistent",
			password: "password123",
			setupMock: func(repo *MockRepository) {
				repo.On("GetUserByUsername", mock.Anything, "nonexistent").Return(nil, nil)
			},
			expectError:   true,
			expectedError: ErrUserNotFound,
		},
		{
			name:     "Incorrect Password",
			username: "testuser",
			password: "wrongpassword",
			setupMock: func(repo *MockRepository) {
				repo.On("GetUserByUsername", mock.Anything, "testuser").Return(mockUser, nil)
			},
			expectError:   true,
			expectedError: ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockRepository)
			tt.setupMock(mockRepo)

			authSvc := NewAuthService(mockRepo, jwtSecret)
			token, err := authSvc.Login(context.Background(), tt.username, tt.password)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
				assert.Empty(t, token)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}
