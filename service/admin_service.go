package service

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/gitslagga/gitbitex-spot/conf"
	"github.com/gitslagga/gitbitex-spot/models"
	"github.com/gitslagga/gitbitex-spot/models/mysql"
	"github.com/pkg/errors"
	"time"
)

func CreateBackendToken(admin *models.Admin) (string, error) {
	claim := jwt.MapClaims{
		"id":       admin.Id,
		"username": admin.Username,
		"password": admin.Password,
		"exp":      time.Now().Add(time.Second * time.Duration(60*60*24*7)).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)
	return token.SignedString([]byte(conf.GetConfig().JwtSecret))
}

func CheckBackendToken(tokenStr string) (*models.Admin, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(conf.GetConfig().JwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claim, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("cannot convert claim to MapClaims")
	}
	if !token.Valid {
		return nil, errors.New("token is invalid")
	}

	usernameValue, found := claim["username"]
	if !found {
		return nil, errors.New("bad token")
	}
	username := usernameValue.(string)

	passwordVal, found := claim["password"]
	if !found {
		return nil, errors.New("bad token")
	}
	password := passwordVal.(string)

	admin, err := GetAdminByUsername(username)
	if err != nil {
		return nil, err
	}
	if admin == nil {
		return nil, errors.New("bad token")
	}
	if admin.Password != password {
		return nil, errors.New("bad token")
	}
	return admin, nil
}

func GetAdminByUsername(username string) (*models.Admin, error) {
	return mysql.SharedStore().GetAdmin(username)
}

func UpdateAdmin(admin *models.Admin) error {
	return mysql.SharedStore().UpdateAdmin(admin)
}
