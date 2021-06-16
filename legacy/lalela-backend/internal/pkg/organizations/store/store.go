package store

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
}

const OrgsServiceProvider = "Forms-Store"

const OrgsCreateOneService = OrgsServiceProvider + ".CreateOne"
const OrgsFindOneService = OrgsServiceProvider + ".FindOne"


type CreateOneRequest struct {
	Name string `json:"name"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}

type FindOneRequest struct {
	Id string `json:"id"`
}

type FindOneResponse struct {
	Name string `json:"name"`
}