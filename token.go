package golib

import (
	"crypto/md5"
	"fmt"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// MyCustomCLaims 载荷
type MyCustomCLaims struct {
	Token string `json:"token"`
	IP    string `json:"ip"`
	jwt.StandardClaims
}

// MyToken .
type MyToken struct{}

// Token .
var Token MyToken

// Encode 获取jwt
func (token *MyToken) Encode(id uint, ip string) (ss string, err error) {
	key := strconv.FormatInt(time.Now().Unix()+int64(id), 10)
	data := []byte(key)
	has := md5.Sum(data)
	t := fmt.Sprintf("%x", has)
	myc := Cache()
	err = myc.Set(t, strconv.FormatInt(int64(id), 10), 3600*2483)
	if err != nil {
		return
	}

	mySigningKey := []byte("GODSON!!!")
	expireTime := time.Now().Add(time.Hour * 24 * 3).Unix()
	claims := MyCustomCLaims{
		t,
		ip,
		jwt.StandardClaims{
			ExpiresAt: expireTime,
			Issuer:    "jwt",
		},
	}
	temp := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err = temp.SignedString(mySigningKey)
	return
}

// Decode 解密jwt
func (token *MyToken) Decode(ss string) (id uint, ip string, err error) {
	temp, err := jwt.ParseWithClaims(ss, &MyCustomCLaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("GODSON!!!"), nil
	})
	if err != nil {
		return 0, "", err
	}
	if claims, ok := temp.Claims.(*MyCustomCLaims); ok && temp.Valid {
		myc := Cache()
		val, err := myc.Get(claims.Token)
		if err != nil {
			return 0, "", err
		}
		tmp, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return 0, "", err
		}
		id = uint(tmp)
		return id, claims.IP, nil
	}
	return 0, "", err
}
