package user

type UseCaseType struct {
	UserName string
	Password string
	Email    string
}

type SignInType struct {
	Email     string
	Password  string
	SessionID string
}
