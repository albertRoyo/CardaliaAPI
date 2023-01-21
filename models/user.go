/*
File		: user.go
Description	: Model file to represent all the user-like objects and their related functions.
It also has implicid functions for the User object.
*/

package models

import (
	"html"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// User DB object
type User struct {
	User_id  uint   `gorm:"primary_key;auto_increment;not_null;" json:"user_id"`
	Username string `gorm:"not_null;unique;" json:"username"`
	Email    string `gorm:"not_null;unique;" json:"email"`
	Password string `gorm:"not_null;" json:"password"`
}

// Used to get the inputs in the frontend
type UserLoginInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Used to get the inputs in the frontend
type UserRegisterInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Used to get the inputs in the frontend
type UserChangePassword struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required"`
}

//////////////////////////////////////////////////////////////////////////////////////////////////////////////////////////

/*
Function	: Save Trade
Description	: Store a new user to the DB.
Self		: User
Parameters 	:
Return     	: User, error
*/
func (u *User) SaveUser() (*User, error) {
	var err error
	err = DB.Create(&u).Error
	if err != nil {
		return &User{}, err
	}
	return u, nil
}

/*
Function	: Before Save
Description	: Actions to do before saving a new user.
Self		: User
Parameters 	:
Return     	: error
*/
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

/*
Function	: Change Password
Description	: Veryfys the old password, encrypts the new one and save the changes.
Self		: User
Parameters 	: old Password, new Password
Return     	: error
*/
func (u *User) ChangePassword(oldPassword string, newPassword string) error {
	// Find user
	DB.First(&u)

	// Verify old password
	if err := VerifyPassword(oldPassword, u.Password); err != nil {
		return err
	}

	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.Password = string(hashedPassword)

	// Save user info
	DB.Save(&u)

	return nil
}

/*
Function	: Verify Password
Description	: Verify the password comparing it to the stored Hashed Password when the user login.
Parameters 	: Password, HasshedPassword
Return     	: error
*/
func VerifyPassword(password, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/*
Function	: Get Username By UserID
Description	: Get the Username of a User giving its UserID.
Parameters 	: UserID
Return     	: Username, error
*/
func GetUsernameByUserID(userID uint) (string, error) {
	var user User
	user.User_id = userID
	err := DB.First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

/*
Function	: Get Email By UserID
Description	: Get the email of a User giving its UserID.
Parameters 	: UserID
Return     	: email, error
*/
func GetEmailByUserID(userID uint) (string, error) {
	var user User
	user.User_id = userID
	err := DB.First(&user).Error
	if err != nil {
		return "", err
	}
	return user.Email, nil
}

/*
Function	: Get UserID By Username
Description	: Get the UserID of a User giving its Username.
Parameters 	: Username
Return     	: UserID, error
*/
func GetUserIDByUsername(username string) (uint, error) {
	user := User{}
	err := DB.Where("username = ?", username).First(&user).Error
	if err != nil {
		return 0, err
	}
	return user.User_id, nil
}
