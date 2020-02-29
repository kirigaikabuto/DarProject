package teachers

import "ldapExample/users"

type Teacher struct {
	users.User
	rating int
}