package domain

type PatchUser struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}