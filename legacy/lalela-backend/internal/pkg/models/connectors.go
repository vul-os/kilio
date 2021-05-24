package models

type ConnectorsModel struct {
	ConnectionId int `json:"connectionId,omitempty"`
	DriverType string `json:"driverType"`
	DriverName string `json:"driverName"`
	ConnectionString string `json:"connectionString"`
}

