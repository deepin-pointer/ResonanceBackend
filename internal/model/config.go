package model

import (
	"github.com/spf13/viper"
)

func InitConfig() {
	viper.SetDefault("bind", ":8000")
	viper.SetDefault("sign_key", "LvE5jSCmICvruafYtFjqE2tIzr3vAk2U")
	viper.SetDefault("static_file", "static.json")
	viper.SetDefault("dynamic_file", "dynamic.bin")
	viper.SetDefault("users", map[string]string{})
}
