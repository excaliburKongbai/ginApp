# 数据库迁移

## 用户表创建

```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,
    nickname VARCHAR(50),
    email VARCHAR(100) UNIQUE,
    mobile VARCHAR(20),
    avatar VARCHAR(255),
    status SMALLINT DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- 创建索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);

-- 添加注释
COMMENT ON TABLE users IS '用户表';
COMMENT ON COLUMN users.id IS '主键ID';
COMMENT ON COLUMN users.username IS '用户名';
COMMENT ON COLUMN users.password IS '密码（加密）';
COMMENT ON COLUMN users.nickname IS '昵称';
COMMENT ON COLUMN users.email IS '邮箱';
COMMENT ON COLUMN users.mobile IS '手机号';
COMMENT ON COLUMN users.avatar IS '头像URL';
COMMENT ON COLUMN users.status IS '状态：1启用 0禁用';
```

## 使用 GORM 自动迁移

在 `cmd/server/main.go` 中添加：

```go
import "ginApp/internal/model"

// 自动迁移
func autoMigrate(db *gorm.DB) {
    if err := db.AutoMigrate(&model.User{}); err != nil {
        log.Fatalf("数据库迁移失败: %v", err)
    }
}
```
