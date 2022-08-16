package model

import "fmt"

type GormConfig struct {
	Debug        bool
	DSN          string
	MaxLifetime  int
	MaxOpenConns int
	MaxIdleConns int
	//TablePrefix  string
}

// SysConfig 读取config.yaml
type SysConfig struct {
	Service    SysService   `yaml:"service" mapstructure:"service"`
	Log        Logger       `yaml:"log" mapstructure:"log"`
	Gorm       GormService  `yaml:"gorm" mapstructure:"gorm"`
	MySql      MySQL        `yaml:"mysql" mapstructure:"mysql"`
	Opa        OpaService   `yaml:"opa" mapstructure:"opa"`
	Mongo      MongoService `yaml:"mongo" mapstructure:"mongo"`
	LdapClient LdapClient   `yaml:"ldapClient" mapstructure:"ldapClient"`

}

type SysService struct {
	Runmode      string `yaml:"runmode" mapstructure:"runmode"`
	Addr         string `yaml:"addr" mapstructure:"addr"`
	Name         string `yaml:"name" mapstructure:"name"`
	Url          string `yaml:"url" mapstructure:"url"`
	MaxPingCount int    `yaml:"maxPingCount" mapstructure:"maxPingCount"`
	JwtSecret    string `yaml:"jwtSecret" mapstructure:"jwtSecret"`
}

// GormService Gorm配置信息
type GormService struct {
	Debug             bool   `yaml:"debug" mapstructure:"debug"`
	DBType            string `yaml:"dbType" mapstructure:"dbType"`
	MaxLifetime       int    `yaml:"maxLifetime" mapstructure:"maxLifetime"`
	MaxOpenConns      int    `yaml:"maxOpenConns" mapstructure:"maxOpenConns"`
	MaxIdleConns      int    `yaml:"maxIdleConns" mapstructure:"maxIdleConns"`
	TablePrefix       string `yaml:"tablePrefix" mapstructure:"tablePrefix"`
	EnableAutoMigrate bool   `yaml:"enableAutoMigrate" mapstructure:"enableAutoMigrate"`
}

// MySQL 读取MySQL部署的配置信息
type MySQL struct {
	Host       string `yaml:"host" mapstructure:"host"`
	Port       int    `yaml:"port" mapstructure:"port"`
	User       string `yaml:"user" mapstructure:"user"`
	Password   string `yaml:"password" mapstructure:"password"`
	DBName     string `yaml:"dbName" mapstructure:"dbName"`
	Parameters string `yaml:"parameters" mapstructure:"parameters"`
}

type Logger struct {
	Dir    string `yaml:"dir" mapstructure:"dir"`
	Remain int    `yaml:"remain" mapstructure:"remain"`
}

// OpaService 读取OPA部署的配置信息
type OpaService struct {
	WatchDirectory string `yaml:"watchDirectory" mapstructure:"watchDirectory"`
}

type MongoService struct {
	Address string `yaml:"address" mapstructure:"address"`
}

// DSN 数据库连接串
func (a MySQL) DSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?%s",
		a.User, a.Password, a.Host, a.Port, a.DBName, a.Parameters)
}

//LdapClient configuration
type LdapClient struct {
	Base             string   `yaml:"base" mapstructure:"base"`
	Host             string   `yaml:"host" mapstructure:"host"`
	Port             int      `yaml:"port" mapstructure:"port"`
	UseSSL           bool     `yaml:"useSsl" mapstructure:"useSsl"`
	SkipTLS          bool     `yaml:"skipTls" mapstructure:"skipTls"`
	BindDN           string   `yaml:"bindDn" mapstructure:"bindDn"`
	BindPassword     string   `yaml:"bindPassword" mapstructure:"bindPassword"`
	UserFilter       string   `yaml:"userFilter" mapstructure:"userFilter"`
	GroupFilter      string   `yaml:"groupFilter" mapstructure:"groupFilter"`
	Attributes       []string `yaml:"attributes" mapstructure:"attributes"`
	ServerName       string   `yaml:"serverName" mapstructure:"serverName"`
	SyncInitPassword string   `yaml:"syncInitPassword" mapstructure:"syncInitPassword"`
	SyncInitRoleID   string   `yaml:"syncInitRoleID" mapstructure:"syncInitRoleID"`
}
