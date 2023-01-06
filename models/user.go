package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	User_id  uint   `gorm:"primary_key;auto_increment;not_null;" json:"user_id"`
	Username string `gorm:"not_null;unique;" json:"username"`
	//Email	string	`gorm:"not_null;unique;" json:"email"`
	Password string `gorm:"not_null;" json:"password"`
}

func GetUsernameByUserID(userID uint) (string, error) {
	var user User
	user.User_id = userID
	err := DB.First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

func GetUserIDByUsername(username string) (uint, error) {
	user := User{}
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.User_id, nil
}

func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (u *User) PrepareGive() {
	u.Password = "******************"
}

func (u *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

func (u *User) BeforeSave() error {

	// Hash Password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	// Remove possible spaces in Username
	u.Username = html.EscapeString(strings.TrimSpace(u.Username))
	return nil

}
