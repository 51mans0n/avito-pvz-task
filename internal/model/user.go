package model

type User struct {
	ID       string `db:"id"`
	Email    string `db:"email"`
	PassHash string `db:"pass_hash"`
	Role     string `db:"role"`
}
