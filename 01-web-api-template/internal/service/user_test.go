package service

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go-practical-roadmap/01-web-api-template/internal/api/dto"
	"go-practical-roadmap/01-web-api-template/internal/model"
)

// MockUserRepository 模拟用户仓库
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id uint) (*model.User, error) {
	args := m.Called(id)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.(*model.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(limit, offset int) ([]model.User, error) {
	args := m.Called(limit, offset)
	result := args.Get(0)
	if result == nil {
		return nil, args.Error(1)
	}
	return result.([]model.User), args.Error(1)
}

func TestUserService_Register_Success(t *testing.T) {
	// 准备测试数据
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	req := &dto.RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "password123",
	}

	// 设置模拟行为
	mockRepo.On("GetByUsername", "testuser").Return(nil, errors.New("user not found"))
	mockRepo.On("GetByEmail", "test@example.com").Return(nil, errors.New("user not found"))
	mockRepo.On("Create", mock.AnythingOfType("*model.User")).Return(nil)

	// 执行测试
	result, err := userService.Register(req)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

func TestUserService_Register_UsernameExists(t *testing.T) {
	// 准备测试数据
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	req := &dto.RegisterRequest{
		Username: "existinguser",
		Email:    "new@example.com",
		Password: "password123",
	}

	existingUser := &model.User{
		ID:       1,
		Username: "existinguser",
		Email:    "existing@example.com",
	}

	// 设置模拟行为
	mockRepo.On("GetByUsername", "existinguser").Return(existingUser, nil)

	// 执行测试
	result, err := userService.Register(req)

	// 验证结果
	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, "username already exists", err.Error())

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}

func TestUserService_GetUserByID_Success(t *testing.T) {
	// 准备测试数据
	mockRepo := new(MockUserRepository)
	userService := NewUserService(mockRepo)

	expectedUser := &model.User{
		ID:       1,
		Username: "testuser",
		Email:    "test@example.com",
	}

	// 设置模拟行为
	mockRepo.On("GetByID", uint(1)).Return(expectedUser, nil)

	// 执行测试
	result, err := userService.GetUserByID(1)

	// 验证结果
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, uint(1), result.ID)
	assert.Equal(t, "testuser", result.Username)
	assert.Equal(t, "test@example.com", result.Email)

	// 验证模拟调用
	mockRepo.AssertExpectations(t)
}