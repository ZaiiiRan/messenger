package chatMember

type role struct {
	Owner     int
	Admin     int
	Member    int
	NotMember int
}

var Roles = role{
	Owner:     3,
	Admin:     2,
	Member:    1,
	NotMember: 0,
}

func GetRoleString(role int) string {
	if role == Roles.Owner {
		return "owner"
	} else if role == Roles.Admin {
		return "admin"
	} else if role == Roles.Member {
		return "member"
	} else {
		return ""
	}
}

func GetRoleValue(role *string) int {
	if *role == "owner" {
		return Roles.Owner
	} else if *role == "admin" {
		return Roles.Admin
	} else if *role == "member" {
		return Roles.Member
	} else {
		return Roles.NotMember
	}
}
