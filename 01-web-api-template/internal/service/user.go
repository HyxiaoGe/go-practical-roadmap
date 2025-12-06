package service

import (
	"errors"

	"go-practical-roadmap/01-web-api-template/internal/model"
	"go-practical-roadmap/01-web-api-template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Register(username, email, password string) (*model.User, error)
	Login(username, password string) (string, error)
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务实例
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// Register 用户注册
func (s *userService) Register(username, email, password string) (*model.User, error) {
	// 检查用户是否已存在
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return nil, errors.New("username already exists")
	}

	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
	}

	// 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 清除密码字段
	user.Password = ""

	return user, nil
}

// Login 用户登录
func (s *userService) Login(username, password string) (string, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// 生成JWT令牌
	// 注意：这里需要导入middleware包来生成令牌
	token := "mock-token" // 实际实现中应调用JWT生成函数

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*model.User, error) {
	return s.userRepo.GetByID(id)
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*model.User, error) {
	return s.userRepo.GetByUsername(username)
}