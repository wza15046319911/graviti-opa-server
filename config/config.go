package config

import (
	"fmt"
	remote "github.com/shima-park/agollo/viper-remote"
	"github.com/spf13/viper"
	"gopa/model"
	CONST "gopa/pkg/constvar"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var Conf *model.SysConfig

func GetConfig() *model.SysConfig {
	return Conf
}

func SetConfig() error {
	if os.Getenv("APP_ENV") == CONST.ConfigOnline {
		fmt.Println("Fetch config from apollo online.")
		if err := FetchApolloConfig(
			os.Getenv("APOLLO_ADDR"),
			os.Getenv("APOLLO_APPID"),
			os.Getenv("APOLLO_NAMESPACE")+"."+CONST.ApolloConfigType,
			CONST.ApolloConfigType); err != nil {
			return err
		}
	} else {
		fmt.Println("fetch config from local config file.")
		if err := FetchLocalConfig(*cfg); err != nil {
			return err
		}
	}
	return nil
}

func FetchLocalConfig(path string) error{
	absPath, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	fmt.Println("dir: ", absPath)
	fmt.Println("Base: ", filepath.Base(path))
	var config model.SysConfig
	f, err := os.Open(absPath)
	if err := yaml.NewDecoder(f).Decode(&config); err != nil {
		return err
	}
	Conf = &config
	fmt.Println(Conf)
	return nil
}

func FetchApolloConfig(addr string, appID string, namespace string, configType string) error {
	var config model.SysConfig
	if err := initApolloConfig(&config, addr, appID, namespace, configType); err != nil {
		return err
	}
	Conf = &config
	fmt.Println("Configuration is: ", Conf.LdapClient.Host)
	return nil
}

// init apollo config
func initApolloConfig(cfg interface{}, addr string, appID string, namespace string, configType string) error {
	remote.SetAppID(appID)
	v := viper.New()
	v.SetConfigType(configType)

	if err := v.AddRemoteProvider("apollo", addr, namespace); err != nil {
		fmt.Println(err)
		return err
	}
	// error handle...
	if err := v.ReadRemoteConfig(); err != nil {
		fmt.Println(err)
		return err
	}

	if err := v.Unmarshal(&cfg); err != nil {
		// panic(err)
		fmt.Println(err)
		return err
	}
	return nil
}
