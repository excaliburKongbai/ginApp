package user

import (
	"context"
	"errors"
	userRequest "ginApp/internal/Dto/Request/user"
	userResponse "ginApp/internal/Dto/Response/user"
	"ginApp/internal/model/user"
	userRepo "ginApp/internal/repository/user"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Service 用户服务接口
type Service interface {
	Login(ctx context.Context, req *userRequest.LoginRequest) (*userResponse.UserResponse, string, error)
	Register(ctx context.Context, req *userRequest.RegisterRequest) (*userResponse.UserResponse, error)
	GetUserInfo(ctx context.Context, userID uint) (*userResponse.UserResponse, error)
	ResetPassword(ctx context.Context, req *userRequest.ResetPasswordRequest) error
}

// service 用户服务实现
type service struct {
	repo userRepo.Repository
}

// NewService 创建用户服务
func NewService(repo userRepo.Repository) Service {
	return &service{repo: repo}
}

// Login 用户登录
func (s *service) Login(ctx context.Context, req *userRequest.LoginRequest) (*userResponse.UserResponse, string, error) {
	// 查找用户
	user, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, "", errors.New("用户名或密码错误")
		}
		return nil, "", err
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, "", errors.New("用户名或密码错误")
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, "", errors.New("账号已被禁用")
	}

	// 生成 token（这里简化处理，实际应该用 JWT）
	token := "mock_token_" + user.Username

	return user.ToResponse(), token, nil
}

// Register 用户注册
func (s *service) Register(ctx context.Context, req *userRequest.RegisterRequest) (*userResponse.UserResponse, error) {
	// 检查用户名是否存在
	existUser, err := s.repo.FindByUsername(ctx, req.Username)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否存在
	if req.Email != "" {
		existEmail, err := s.repo.FindByEmail(ctx, req.Email)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existEmail != nil {
			return nil, errors.New("邮箱已被注册")
		}
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &user.User{
		Username: req.Username,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
		Email:    req.Email,
		Mobile:   req.Mobile,
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

// GetUserInfo 获取用户信息
func (s *service) GetUserInfo(ctx context.Context, userID uint) (*userResponse.UserResponse, error) {
	user, err := s.repo.FindByID(ctx, userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("用户不存在")
		}
		return nil, err
	}

	return user.ToResponse(), nil
}

// ResetPassword 重置密码
func (s *service) ResetPassword(ctx context.Context, req *userRequest.ResetPasswordRequest) error {
	// 查找用户
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("邮箱未注册")
		}
		return err
	}

	// TODO: 验证验证码（这里简化处理）
	if req.Code != "123456" {
		return errors.New("验证码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	return s.repo.Update(ctx, user)
}
