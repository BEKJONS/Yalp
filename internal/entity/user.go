package entity

// Define user role, user type, and user status as constants for type safety
const (
	UserRoleUser       = "user"
	UserRoleAdmin      = "admin"
	UserRoleSuperAdmin = "superadmin"

	UserTypeUser  = "user"
	UserTypeAdmin = "admin"

	UserStatusActive   = "active"
	UserStatusBlocked  = "blocked"
	UserStatusInVerify = "inverify"
)

type User struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	FullName    string `json:"full_name"`
	UserType    string `json:"user_type"`
	UserRole    string `json:"user_role"`
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
	DeletedAt   string `json:"deleted_at,omitempty"` // can be null
}

type UserSingleRequest struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

type UserList struct {
	Items []User `json:"users"`
	Count int    `json:"count"`
}
