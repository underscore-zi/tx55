package admin

type Privilege string

const (
	PrivNone           Privilege = ""
	PrivAll            Privilege = "all_privileges"
	PrivUpdateProfiles Privilege = "update_profiles"
	PrivFullIPs        Privilege = "full_ips"
	PrivSearchByIP     Privilege = "search_by_ip"
	PrivReadBans       Privilege = "read_bans"
	PrivUpdateBans     Privilege = "update_bans"
)

func (u *User) HasPrivilege(p Privilege) bool {
	if u.Role.AllPrivileges {
		return true
	}

	switch p {
	case PrivNone:
		return true
	case PrivAll:
		return u.Role.AllPrivileges
	case PrivUpdateProfiles:
		return u.Role.UpdateProfiles
	case PrivReadBans:
		return u.Role.ReadBans
	case PrivUpdateBans:
		return u.Role.UpdateBans
	case PrivFullIPs:
		return u.Role.FullIPs
	case PrivSearchByIP:
		return u.Role.SearchByIP
	}

	return false
}
