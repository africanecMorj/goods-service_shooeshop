package domain

type ResetRequest struct {
    Email string `json:"email"`
}

type ResetPassword struct {
    Email    string `json:"email"`
    Password string `json:"password"`
}