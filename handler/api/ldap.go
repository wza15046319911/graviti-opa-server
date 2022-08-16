package api

import (
	"fmt"
	ldap "github.com/jtblin/go-ldap-client"
	"gopa/config"
)

//LdapAuth 加载LDAP配置参数
func LdapAuth(username, password string) error {
	ldapClient := &ldap.LDAPClient{
		Base:         config.GetConfig().LdapClient.Base,
		Host:         config.GetConfig().LdapClient.Host,
		Port:         config.GetConfig().LdapClient.Port,
		UseSSL:       config.GetConfig().LdapClient.UseSSL,
		SkipTLS:      config.GetConfig().LdapClient.SkipTLS,
		BindDN:       config.GetConfig().LdapClient.BindDN,
		BindPassword: config.GetConfig().LdapClient.BindPassword,
		UserFilter:   config.GetConfig().LdapClient.UserFilter,
		GroupFilter:  config.GetConfig().LdapClient.GroupFilter,
		Attributes:   config.GetConfig().LdapClient.Attributes,
		ServerName:   config.GetConfig().LdapClient.ServerName,
	}
	defer ldapClient.Close()
	_, _, err := ldapClient.Authenticate(username, password)
	fmt.Println(err)
	if err != nil {
		return err
	}

	return nil
}
