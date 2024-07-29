package model

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
)

type User struct {
	ID          int64
	Name        string
	Birth       string
	PhoneNumber string
	Password    string
	IsDeleted   bool
}

func (user User) IsCorrectPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

/*
일반적인 휴대전화 5,6,9,10번쨰 마킹 0107**11**4
*/
func (user User) MaskPhoneNumber() string {
	parts := strings.Split(user.PhoneNumber, "")
	for _, idx := range []int{4, 5, 8, 9} {
		parts[idx] = "*"
	}
	return strings.Join(parts, "")
}
