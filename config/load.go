package config

import "github.com/spf13/pflag"

var (
	cfg = pflag.StringP("config", "c", "./conf/config.yaml", "config file path.")
)
