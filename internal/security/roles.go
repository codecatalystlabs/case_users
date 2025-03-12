package security

func GetRoles(user int, role string) int {
	if user > 0 && role != "" {
		return 1
	} else {
		return 0
	}
}
