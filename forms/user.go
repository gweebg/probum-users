package forms

type UserSignup struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required"`
	UId      string `json:"uid" binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role" binding:"required"`
}

type UserUpdate struct {
	Name     *string `json:"name,omitempty"`
	Email    *string `json:"email,omitempty"`
	Password *string `json:"password,omitempty"`
}

type AuthUser struct {
	UId      string
	Password string
}
