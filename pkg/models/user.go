package models

type User struct {
	Id        int
	FirstName string
	LastName  string
	Email     string
}

func (u *User) GetFullName() string {
	return u.FirstName + " " + u.LastName
}
