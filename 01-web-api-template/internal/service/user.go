package service

import (
	"errors"

	"go-practical-roadmap/01-web-api-template/internal/api/dto"
	"go-practical-roadmap/01-web-api-template/internal/middleware"
	"go-practical-roadmap/01-web-api-template/internal/model"
	"go-practical-roadmap/01-web-api-template/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Register(req *dto.RegisterRequest) (*dto.UserProfileResponse, error)
	Login(req *dto.LoginRequest) (string, error)
	GetUserByID(id uint) (*dto.UserProfileResponse, error)
	GetUserByUsername(username string) (*dto.UserProfileResponse, error)
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
func (s *userService) Register(req *dto.RegisterRequest) (*dto.UserProfileResponse, error) {
	// 检查用户是否已存在
	if _, err := s.userRepo.GetByUsername(req.Username); err == nil {
		return nil, errors.New("username already exists")
	}

	if _, err := s.userRepo.GetByEmail(req.Email); err == nil {
		return nil, errors.New("email already exists")
	}

	// 密码加密
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	// 保存到数据库
	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	// 返回用户信息（不包含密码）
	response := &dto.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}

// Login 用户登录
func (s *userService) Login(req *dto.LoginRequest) (string, error) {
	// 获取用户
	user, err := s.userRepo.GetByUsername(req.Username)
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return "", errors.New("invalid username or password")
	}

	// 生成JWT令牌
	token, err := middleware.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(id uint) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	response := &dto.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}

// GetUserByUsername 根据用户名获取用户
func (s *userService) GetUserByUsername(username string) (*dto.UserProfileResponse, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, err
	}

	response := &dto.UserProfileResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
	}

	return response, nil
}