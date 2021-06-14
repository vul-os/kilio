package controllers

import (
	"lalela-backend/internal/pkg/models"
	"lalela-backend/internal/pkg/utils"
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

// Define AR Structs
type DashboardsAR struct {
	UserEmail  string             `json:"user_email"`
	Dashboards []models.Dashboard `json:"dashboards"`
}

type UsersAR struct {
	UserEmail string        `json:"user_email"`
	Users     []models.User `json:"users"`
}

type PermissionsAR struct {
	Permissions []models.DashboardPermission `json:"permissions"`
}

func CheckJson(w http.ResponseWriter, r *http.Request, i interface{}) (interface{}, bool) {
	// Define Struct Map
	err := json.NewDecoder(r.Body).Decode(&i)

	if err != nil {
		resp, err := utils.GetError(400, "en")
		if err != nil {
			json.NewEncoder(w).Encode(err)
			return i, true
		}
		json.NewEncoder(w).Encode(resp)
		return i, true
	}
	return i, false
}

func GetEnvVar(t string) string {
	return viper.Get(t).(string)
}
