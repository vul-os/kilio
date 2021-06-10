package store

type Store interface {
	CreateOne(CreateOneRequest) (*CreateOneResponse, error)
	GetOne(FindOneRequest) (*FindOneResponse, error)
}

const FormsServiceProvider = "Forms-Store"

const FormsCreateOneService = FormsServiceProvider + ".CreateOne"
const FormsFindOneService = FormsServiceProvider + ".FindOne"


type CreateOneRequest struct {
	Name string `json:"name"`
	Scheme interface{} `json:"scheme"`
	UiScheme interface{} `json:"ui_scheme"`
}

type CreateOneResponse struct {
	Id string `json:"id"`
}

type FindOneRequest struct {
	Id string `json:"id"`
}

type FindOneResponse struct {
	Name string `json:"name"`
	Scheme interface{} `json:"scheme"`
	UiScheme interface{} `json:"ui_scheme"`
}