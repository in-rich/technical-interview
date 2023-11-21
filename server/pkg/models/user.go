package models

type User struct {
	ID       string `json:"id" firestore:"id"`
	Email    string `json:"email" firestore:"email"`
	Username string `json:"username" firestore:"username"`
	Password string `json:"-" firestore:"password"`
}
