package templates

func doesStrMatchAny(userRole string, allowedRoles ...string) bool {
	var	role	string

	for _, role = range allowedRoles {
		if userRole == role {
			return true
		}
	}
	return false
}
