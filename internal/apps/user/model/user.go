package model

import (
	_const "Aurora/internal/pkg/const"
	"errors"
	"gorm.io/gorm"
	"time"
)

type UserDal struct {
	conn *gorm.DB
}

type User struct {
	ID        string `gorm:"primary_key;column:id"`
	Account   string `gorm:"type:varchar(20); unique_index; column:account; default:'' "`
	Username  string `gorm:"type:varchar(20); unique_index; column:username; default:'默认名' " `
	Password  string `gorm:"type:varchar(255); column:password"`
	Avatar    string `gorm:"default:''; column:avatar"`
	Gender    uint8  `gorm:"type:tinyint(3);default:0; column:gender"`
	Mobile    string `gorm:"type:varchar(20);default:'';unique_index;column:mobile"`
	Email     string `gorm:"type:varchar(30);default:'';unique_index;column:email"`
	Introduce string `gorm:"default:'';column:introduce"`
	Status    uint8  `gorm:"type:tinyint(3);default:1;column:status"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type FilterUser struct {
	ID       []string
	Account  []string
	Username []string
	Gender   []uint8
	Mobile   []string
	Email    []string
	Paging
}

type Paging struct {
	PageSize int `default:"1"`
	PageNum  int `default:"20"`
}

type UserPaging struct {
	User
	Paging
}

func NewUserDal(conn *gorm.DB) *UserDal {
	return &UserDal{
		conn: conn,
	}
}

func (u *User) TableName() string {
	return "user"
}

func (u *UserDal) InsertUser(user *User) (*User, error) {
	if err := u.conn.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDal) GetUserById(id string, page *Paging) (*User, error) {
	user := &User{}
	limit, offset := VerifyPage(page)
	err := u.conn.Where("id = ? and status = ?", id, _const.StatusUse).First(user).Limit(limit).Offset(offset).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDal) UpdateUser(user *User) error {
	err := u.conn.Updates(user).Error
	if err != nil {
		return err
	}
	return nil
}

// DeleteUserByID set status as 1
func (u *UserDal) DeleteUserByID(id string) error {
	err := u.conn.Model(&User{}).Update("status", _const.StatusDelete).Where("id = ?", id).Error
	if err != nil {
		return err
	}
	return nil
}

func (u *UserDal) GetCountByAccount(account string) (int64, error) {
	var total int64
	err := u.conn.Model(&User{}).Count(&total).Error
	if err != nil {
		return -1, err
	}
	return total, nil
}

func (u *UserDal) GetUserInfoByAccount(account string) (*User, error) {
	user := &User{}
	err := u.conn.Where("account = ? and status = ?", account, _const.StatusUse).First(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDal) BatchGetUserInfoById(ids []string) ([]User, error) {
	user := make([]User, 0)
	err := u.conn.Where("id in ? and status = ?", ids, _const.StatusDelete).Find(user).Error
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (u *UserDal) FilterGetUserInfo(filter *FilterUser) (int64, []User, error) {
	if filter == nil {
		return 0, nil, errors.New("empty filter")
	}
	db := u.conn.Model(&User{}).Where("status = ?", _const.StatusUse)
	if len(filter.ID) > 0 {
		db = db.Where("id in ( ? )", filter.ID)
	}
	if len(filter.Email) > 0 {
		db = db.Where("email in ( ? ) ", filter.Email)
	}
	if len(filter.Mobile) > 0 {
		db = db.Where("mobile in ( ? )", filter.Mobile)
	}
	if len(filter.Username) > 0 {
		for _, name := range filter.Username {
			db = db.Or("username LIKE ?", "%"+name+"%")
		}
	}

	if len(filter.Account) > 0 {
		db = db.Where("account in ( ? )", filter.Account)
	}

	limit, offset := VerifyPage(&Paging{filter.PageSize, filter.PageNum})
	db = db.Limit(limit).Offset(offset)

	var users []User
	var total int64
	err := db.Order("created_at DESC").Find(&users).
		Offset(0).Limit(-1).Count(&total).Error
	return total, users, err
}

func VerifyPage(page *Paging) (limit int, offset int) {
	if page == nil {
		return 20, 0
	}

	if page.PageNum <= 0 {
		page.PageNum = 1
	}

	if page.PageSize <= 0 {
		page.PageSize = 20
	}

	return page.PageSize, (page.PageNum - 1) * page.PageSize
}
