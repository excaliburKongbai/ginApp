package user

import (
	"context"
	"ginApp/internal/model/user"
	"gorm.io/gorm"
)

// Repository 用户仓储接口
type Repository interface {
	Create(ctx context.Context, user *user.User) error
	FindByID(ctx context.Context, id uint) (*user.User, error)
	FindByUsername(ctx context.Context, username string) (*user.User, error)
	FindByEmail(ctx context.Context, email string) (*user.User, error)
	Update(ctx context.Context, user *user.User) error
	Delete(ctx context.Context, id uint) error
}

// repository 用户仓储实现
type repository struct {
	db *gorm.DB
}

// NewRepository 创建用户仓储
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

// Create 创建用户
func (r *repository) Create(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

// FindByID 根据 ID 查找用户
func (r *repository) FindByID(ctx context.Context, id uint) (*user.User, error) {
	var user user.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername 根据用户名查找用户
func (r *repository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var user user.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail 根据邮箱查找用户
func (r *repository) FindByEmail(ctx context.Context, email string) (*user.User, error) {
	var user user.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update 更新用户
func (r *repository) Update(ctx context.Context, user *user.User) error {
	return r.db.WithContext(ctx).Save(user).Error
}

// Delete 删除用户（软删除）
func (r *repository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&user.User{}, id).Error
}
