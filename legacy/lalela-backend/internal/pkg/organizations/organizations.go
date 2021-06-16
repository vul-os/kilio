package organizations

type Organizations struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"organization_name"`
}
