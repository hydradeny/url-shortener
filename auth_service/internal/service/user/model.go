package user

type User struct {
	ID    uint32
	Email string
}

type RawUser struct {
	ID       uint32
	Email    string
	PassHash []byte
}

type CreateUser struct {
	Email    string
	Password string
}

type CheckPassword struct {
	Email    string
	Password string
}
