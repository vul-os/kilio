package forms


type Forms struct {
	FormName  string             `json:"form_name" bson:"form_name"`
	Scheme    interface{}        `json:"scheme" bson:"scheme"`
	UiScheme  interface{}        `json:"ui_scheme" bson:"ui_scheme"`
}






