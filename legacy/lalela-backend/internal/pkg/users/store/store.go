package store

import (
	"time"
)

type Store interface {
	CreateUser(CreateUserRequest) (*CreateUserResponse, error)
	FindUser(FindUserRequest) (*FindUserResponse, error)
}

const UserServiceProvider = "User-Store"

const UserCreateUserService = UserServiceProvider + ".CreateUser"
const UserFindUserService = UserServiceProvider + ".FindUser"

type CreateUserRequest struct {
	ID              string 				`json:"id"`
	OrganizationId  string 				`json:"organization_id"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
	Password        string             `json:"password"`
	Email           string             `json:"email"`
	RoleID          string             `json:"role_id"`
	ValidationToken string             `json:"validation_token"`
	EmailToken      string             `json:"email_token"`
	RefreshToken    string             `json:"refresh_token"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type CreateUserResponse struct {
	Id 				string 				`json:"id"`
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
}
type FindUserRequest struct {
	Id string `json:"id"`
}

type FindUserResponse struct {
	FirstName       string             `json:"first_name"`
	LastName        string             `json:"last_name"`
}

/*func CreateOne(email string, password string) (string, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	var user users.User
	var collection = database.OpenCollection("user")

	user.ID = primitive.NewObjectID()
	user.Email = email


	token, _ := auth.GenerateToken(email)
	user.ValidationToken = token
	pass := auth.HashPassword(password)
	user.Password = pass
	user.OrganizationId = primitive.NilObjectID

	user.CreatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	user.UpdatedAt, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))

	insertId, err := collection.InsertOne(ctx, user)
	if err != nil {
		cancel()
		return "", err
	}
	defer cancel()

	return cast.ToString(insertId.InsertedID), err
}

func FindUserByEmail(email string) (users.User, error) {
	var collection = database.OpenCollection("user")
	var foundUser users.User
	err := collection.FindOne(context.Background(), bson.M{"email": email}).Decode(&foundUser)
	if err != nil {
		return users.User{}, err
	}
	return foundUser, err
}*/
