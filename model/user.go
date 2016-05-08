package model

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"strings"

	"github.com/mephux/kolide/shared/base"
)

// User database table schema
type User struct {
	Id   int64
	Name string `xorm:"NOT NULL" json:"name"`

	Key         string `xorm:"UNIQUE INDEX NOT NULL" json:"key"`
	Email       string `xorm:"UNIQUE NOT NULL" json:"email"`
	Avatar      string `xorm:"VARCHAR(2048) NOT NULL" json:"avatar"`
	AvatarEmail string `xorm:"NOT NULL" json:"-"`

	Password string `xorm:"NOT NULL" json:"-"`
	Hash     string `xorm:"VARCHAR(10)" json:"-"`
	Salt     string `xorm:"VARCHAR(10)" json:"-"`

	Admin   bool `json:"admin"`
	Enabled bool `json:"enabled"`
}

var (
	// ErrUserNotExist user not found error message
	ErrUserNotExist = errors.New("User does not exist")
)

// AvatarLink Users avatar link. Using gravatar
func (u *User) AvatarLink() string {
	return "//1.gravatar.com/avatar/" + u.Avatar
}

// GetUserById returns the user object by given ID if exists.
func GetUserById(id int64) (*User, error) {
	u := new(User)

	has, err := x.Id(id).Get(u)

	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}

	return u, nil
}

// GetUserByEmail returns the user object by given e-mail if exists.
func GetUserByEmail(email string) (*User, error) {
	if len(email) == 0 {
		return nil, ErrUserNotExist
	}

	user := &User{
		Email: strings.ToLower(email),
	}

	has, err := x.Get(user)
	if err != nil {
		return nil, err
	} else if !has {
		return nil, ErrUserNotExist
	}
	return user, nil
}

// UserLogin validates user name and password.
func UserLogin(email, password string) (*User, error) {
	var u *User

	u = &User{
		Email: email,
	}

	userExists, err := x.Get(u)

	if err != nil {
		return nil, err
	}

	if userExists {
		if u.ValidatePassword(password) {
			return u, nil
		}

		return nil, fmt.Errorf("authentication error")
	}

	return nil, fmt.Errorf("user does not exist")
}

// Exists returns true or false based on the user context
func (u *User) Exists() bool {
	userExists, err := x.Get(u)

	if err != nil {
		return false
	}

	return userExists
}

// ValidatePassword checks if given password matches the one belongs to the user.
func (u *User) ValidatePassword(passwd string) bool {
	newUser := &User{Password: passwd, Salt: u.Salt}
	newUser.EncodePassword()
	return u.Password == newUser.Password
}

// EncodePassword encodes password to safe format.
func (u *User) EncodePassword() {
	newPasswd := base.PBKDF2([]byte(u.Password), []byte(u.Salt), 10000, 50, sha256.New)
	u.Password = fmt.Sprintf("%x", newPasswd)
}

// IsEmailUsed returns true if the e-mail has been used.
func IsEmailUsed(email string) (bool, error) {
	if len(email) == 0 {
		return false, nil
	}

	return x.Get(&User{
		Email: email,
	})
}

// GetUserSalt returns a user salt token
func GetUserSalt() string {
	return base.GetRandomString(10)
}

// CreateUser creates record of a new user.
func CreateUser(u *User) error {
	isExist, err := IsEmailUsed(u.Email)

	if err != nil {
		return err
	} else if isExist {
		return fmt.Errorf("e-mail has been used [email: %s]", u.Email)
	}

	u.Avatar = base.EncodeMd5(u.Email)
	u.AvatarEmail = u.Email
	u.Hash = GetUserSalt()
	u.Salt = GetUserSalt()
	u.EncodePassword()
	u.Key = UniqueKey()

	sess := x.NewSession()
	defer sess.Close()
	if err = sess.Begin(); err != nil {
		return err
	}

	// Auto-set admin for user whose ID is 1.
	if u.Id == 1 {
		u.Admin = true
		u.Enabled = true
		_, err = x.Id(u.Id).UseBool().Update(u)
	}

	if _, err = sess.Insert(u); err != nil {
		sess.Rollback()
		return err
	}

	return sess.Commit()
}
