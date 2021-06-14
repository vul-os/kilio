package store

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	FindOne(FindOneRequest) (*FindOneResponse, error)
}

const FormsServiceProvider = "Submissions-Store"

const FormsCreateOneService = FormsServiceProvider + ".CreateOne"
const FormsFindOneService = FormsServiceProvider + ".FindOne"

type CreateOneRequest struct {
	FormId string      `json:"form_id"`
	OrgId  string      `json:"org_id"`
	Data   interface{} `json:"data"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}

type FindOneRequest struct {
	Id string `json:"id"`
}

type FindOneResponse struct {
	FormId string      `json:"form_id"`
	OrgId  string      `json:"org_id"`
	Data   interface{} `json:"data"`
}
