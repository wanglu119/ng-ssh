package viper

import (
	"github.com/spf13/viper"
)

type Viper = viper.Viper
type ConfigParseError = viper.ConfigParseError

var v *Viper

func init() {
	v = viper.New()
}

func GetViper() *Viper {
	return v
}