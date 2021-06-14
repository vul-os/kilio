package users

type User struct {
	ID              string   `json:"id" bson:"id"`
	Name            string   `json:"name" bson:"name"`
	Email           string   `json:"email" bson:"email"`
	Password        []byte   `json:"-" bson:"password"`
	OrganizationIds []string `json:"organization_ids" bson:"organization_ids"`
	ResetToken      string   `json:"refresh_token" bson:"refresh_token"`
}
