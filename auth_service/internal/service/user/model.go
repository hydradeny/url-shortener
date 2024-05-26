package user

type User struct {
	ID       uint32
	Email    string
	Password string
}

type CreateUser struct {
	Email    string
	Password string
}
