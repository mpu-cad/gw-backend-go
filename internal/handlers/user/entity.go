package user

type loginUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registrationUserRequest struct {
	Name     string `json:"name"`
	Surname  string `json:"surname"`
	LastName string `json:"last_name"`
	Login    string `json:"login"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
}
