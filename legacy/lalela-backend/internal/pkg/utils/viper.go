package utils

import "github.com/spf13/viper"

//TODO: ???
func GetEnvVar(t string) string {
	return viper.Get(t).(string)
}
