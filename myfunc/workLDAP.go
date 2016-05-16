package myfunc

import (
	"fmt"
	"github.com/gogits/gogs/modules/ldap"
)

const constTimeSearch = 60

var Base_dn string = "dc=npp"
var Filter string = "(objectClass=user)"
var IntColObjects int = 1000

func WorkLDAP() (*ldap.Conn, []*ldap.Entry, error) {
	const (
		UserName     = "CN=oitt,OU=oit,DC=npp"
		UserPassword = "05112015"
	)
	var ldap_server string = "ad2.npp"
	var ldap_port uint16 = 389
	var attributes []string = []string{
		"cn",
		"objectCategory",
		"sn",
		"sAMAccountName",
		"memberOf",
		"member"}

	localLDAPconn, ldapErr := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldap_server, ldap_port))

	if ldapErr != nil {
		CLog.PrintLog(ldapErr)
		return nil, nil, ldapErr
	}

	if ldapErr := localLDAPconn.Bind(UserName, UserPassword); ldapErr != nil {
		CLog.PrintLog(ldapErr)
		return nil, nil, ldapErr
	}

	search := ldap.NewSearchRequest(
		Base_dn,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, IntColObjects, constTimeSearch, false, Filter, attributes, nil)

	CLog.PrintLog("search = ", search)

	sr, ldapErr := localLDAPconn.Search(search)

	if sr == nil {
		CLog.PrintLog(ldapErr)
		return localLDAPconn, nil, ldapErr
	}

	if ldapErr != nil {
		CLog.PrintLog(ldapErr)
		return localLDAPconn, sr.Entries, ldapErr
	}

	return localLDAPconn, sr.Entries, nil
}
