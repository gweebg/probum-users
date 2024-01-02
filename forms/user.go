package forms

type UserSignup struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	UId      string `json:"uid"      binding:"required"`
	Password string `json:"password" binding:"required"`
	Role     string `json:"role"     binding:"required,oneof=tech admin lstudent mstudent"`
}

type UserUpdate struct {
	Name     *string `json:"name,omitempty"     binding:"required"`
	Email    *string `json:"email,omitempty"    binding:"required,email"`
	Password *string `json:"password,omitempty" binding:"required"`
}

type UserLogin struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}
