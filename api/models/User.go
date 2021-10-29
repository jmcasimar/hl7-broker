package models

import (
	"errors"
	"html"
	"time"

	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint      `gorm:"primary_key;auto_increment" json:"id"`
	Name      string    `gorm:"size:255;not null" json:"name"`
	Title     string    `gorm:"size:255;" json:"title"`
	NIC       string    `gorm:"size:25;not null;unique;unique_index" json:"nic"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	Password  string    `gorm:"size:100;not null;" json:"-"`
}

func (u *User) Prepare() {
	u.ID = 0
	u.Name = html.EscapeString(u.Name)
	u.NIC = html.EscapeString(u.NIC)
	u.Title = html.EscapeString(u.Title)
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()
}

func (u *User) Validate(action string) error {
	switch action {
	case "update":
		if u.Name == "" {
			return errors.New("name is required")
		}
		if u.Title == "" {
			return errors.New("title is required")
		}
		if u.NIC == "" {
			return errors.New("nic is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
	default:
		if u.Name == "" {
			return errors.New("name is required")
		}
		if u.Title == "" {
			return errors.New("title is required")
		}
		if u.NIC == "" {
			return errors.New("nic is required")
		}
		if u.Password == "" {
			return errors.New("password is required")
		}
	}

	return nil
}

func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) BeforeSave() error {
	hashedPassword, err := Hash(u.Password)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)
	return nil
}

func (u *User) SaveUser(db *gorm.DB) (*User, error) {
	var err error
	err = db.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) FindAllUsers(db *gorm.DB, limit int, offset int) (*[]User, error) {
	var err error
	var users []User
	err = db.Model(&User{}).Limit(limit).Offset(offset).Find(&users).Error
	if err != nil {
		return &[]User{}, err
	}
	return &users, err
}

func (u *User) FindUserByID(db *gorm.DB, uid uint32) (*User, error) {
	var err error
	err = db.Model(User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	if gorm.IsRecordNotFoundError(err) {
		return &User{}, errors.New("user not found")
	}
	return u, err
}

func (u *User) UpdateAUser(db *gorm.DB, uid uint32) (*User, error) {

	// To hash the password
	err := u.BeforeSave()
	if err != nil {
		log.Fatal(err)
	}
	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).UpdateColumns(
		map[string]interface{}{
			"password":  u.Password,
			"nic":       u.NIC,
			"name":      u.Name,
			"title":     u.Title,
			"update_at": time.Now(),
		},
	)
	if db.Error != nil {
		return &User{}, db.Error
	}
	// This is the display the updated user
	err = db.Model(&User{}).Where("id = ?", uid).Take(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) DeleteAUser(db *gorm.DB, uid uint32) (int64, error) {

	db = db.Model(&User{}).Where("id = ?", uid).Take(&User{}).Delete(&User{})

	if db.Error != nil {
		return 0, db.Error
	}
	return db.RowsAffected, nil
}
