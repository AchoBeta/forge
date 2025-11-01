package util

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// 密码加密
// 使用 bcrypt 生成密码哈希
func HashPassword(password string) (string, error) {
	if password == "" {
		return "", errors.New("password 不能为空")
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// 校验明文密码与哈希是否匹配
func ComparePassword(hash string, plain string) (bool, error) {
	if hash == "" || plain == "" {
		return false, errors.New("hash 或 password 不能为空")
	}
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(plain))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

//

// jwt  后续登录给到token
// JWT工具类
type JWTUtil struct {
	secretKey   []byte //密钥
	expireHours int    //过期时间
}

// JWT声明
type Claims struct {
	UserID string `json:"user_id"` //用户唯一标识  解析token识别用户
	jwt.RegisteredClaims
}

// 创建JWT工具实例
func NewJWTUtil(secretKey string, expireHours int) *JWTUtil {
	return &JWTUtil{
		secretKey:   []byte(secretKey),
		expireHours: expireHours,
	}
}

// GenerateToken 生成jwt令牌
func (j *JWTUtil) GenerateToken(userID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID 不能为空")
	}

	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secretKey)
}

/* 验证JWT令牌
func (j *JWTUtil) ValidateToken(tokenString string) (*Claims, error)
*/

/* 刷新令牌
func (j *JWTUtil) RefreshToken(tokenString string) (string, error)
*/
