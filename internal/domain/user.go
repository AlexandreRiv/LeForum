package domain

type User struct {
	ID       int
	Username string
	Email    string
	Password string
	DarkMode bool
	Role     RoleType
}

// RoleType définit un type pour les rôles utilisateur
type RoleType string

const (
	RoleUser      RoleType = "user"
	RoleModerator RoleType = "moderator"
	RoleAdmin     RoleType = "admin"
)
