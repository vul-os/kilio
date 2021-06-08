package auth

import (
	"context"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

// SignedDetails
type SignedDetails struct {
	Email string
	Id    string
	jwt.StandardClaims
}

var SecretKey string

// GenerateAllTokens generates both teh detailed token and refresh token
func GenerateToken(email string) (signedToken string, err error) {
	claims := &SignedDetails{
		Email: email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SecretKey))

	if err != nil {
		log.Panic(err)
		return
	}

	return token, err
}

func ValidateJWTRequest(r *http.Request) *SignedDetails {
	token := getTokenFromRequest(r)
	return ValidateToken(token)
}

func getTokenFromRequest(r *http.Request) string {
	const BEARER_SCHEMA = "Bearer "
	authHeader := r.Header.Get("Authorization")
	return authHeader[len(BEARER_SCHEMA):]
}

//ValidateToken validates the jwt token
func ValidateToken(signedToken string) (claims *SignedDetails) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
		},
	)

	if err != nil {
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		fmt.Println("the token is invalid")
		return
	}

	if claims.ExpiresAt < time.Now().Local().Unix() {
		fmt.Println("token is expired")
		return
	}

	return claims
}

func UpdateToken(collection *mongo.Collection, signedToken string, id string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

	var updateObj primitive.D
	updateObj = append(updateObj, bson.E{Key: "ValidationToken", Value: signedToken})
	updateObj = append(updateObj,
		bson.E{Key: "updated_at", Value: time.Now().Format(time.RFC3339)},
	)

	upsert := true
	filter := bson.M{"_id": id}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := collection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{"$set", updateObj},
		},
		&opt,
	)
	defer cancel()

	if err != nil {
		log.Panic(err)
		return
	}

	return
}
