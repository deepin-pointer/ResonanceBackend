package model

import (
	"crypto/rand"
	"encoding/base64"
	"math"

	"github.com/spf13/viper"
)

func randomBase64String(l int) string {
	buff := make([]byte, int(math.Ceil(float64(l)/float64(1.33333333333))))
	rand.Read(buff)
	str := base64.RawURLEncoding.EncodeToString(buff)
	return str[:l]
}

func InitConfig() {
	viper.SetDefault("bind", ":8000")
	viper.SetDefault("sign_key", randomBase64String(32))
	viper.SetDefault("static_file", "static.json")
	viper.SetDefault("dynamic_file", "dynamic.bin")
	viper.SetDefault("users", map[string]string{})
}
