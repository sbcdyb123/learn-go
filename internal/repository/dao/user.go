package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserDuplicateEmail = errors.New("邮箱冲突")
	ErrUserNotFound       = gorm.ErrRecordNotFound
)

type UserDao struct {
	db *gorm.DB
}

func NewUserDao(db *gorm.DB) *UserDao {
	return &UserDao{
		db: db,
	}
}
func (dao *UserDao) Inset(ctx context.Context, u User) error {
	now := time.Now().UnixMilli()
	u.UTime = now
	u.Ctime = now
	err := dao.db.WithContext(ctx).Create(&u).Error
	var mysqlErr *mysql.MySQLError
	if errors.As(err, &mysqlErr) {
		const uniqueConflictsErrNo uint16 = 1062
		if mysqlErr.Number == uniqueConflictsErrNo {
			return ErrUserDuplicateEmail
		}
		return err
	}
	return nil
}

func (dao *UserDao) FindByEmail(ctx context.Context, email string) (User, error) {
	var u User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&u).Error
	return u, err
}

func (dao *UserDao) Updates(ctx context.Context, u User) error {
	err := dao.db.WithContext(ctx).Where("id = ?", u.Id).Updates(u).Error
	return err
}

type User struct {
	Id       int64  `gorm:"primaryKey,autoIncrement"`
	Email    string `gorm:"unique"`
	Password string
	Username string
	BirthDay int64
	Intro    string
	Ctime    int64
	UTime    int64
}
