package utils
//
// import (
// 	"time"
// 	"gitlab.com/paiduay/queq-hospital-api/config"
// 	"gitlab.com/paiduay/queq-hospital-api/models"
//
// 	"github.com/dgrijalva/jwt-go"
// )
//
// // JwtCustomClaims - struct for token parsing.
// type JwtCustomClaims struct {
// 	jwt.StandardClaims
// 	ID          uint32 `json:"userID"`
// 	Username    string `json:"username"`
// 	UserTypeID uint8 `json:"userTypeID"`
// }
//
// // IssueJWT - create token when a user login
// func IssueJWT(user *models.User) (string, error) {
// 	tokenDuration, _ := time.ParseDuration(config.Configs.JWT.Lifetime)
//
// 	claim := JwtCustomClaims{
// 		jwt.StandardClaims{
// 			ExpiresAt: time.Now().Add(tokenDuration).Unix(),
// 		},
// 		user.ID,
// 		*user.Username,
// 		user.UserTypeID,
// 	}
// 	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claim)
// 	signedToken, err := token.SignedString([]byte(config.Configs.JWT.Secret))
// 	if err != nil {
// 		return "", err
// 	}
// 	return signedToken, nil
// }
