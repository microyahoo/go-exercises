package main

import (
	"fmt"
)

type Notifier interface {
	Notify() error
	Print() string
}

type User struct {
	Name  string
	Email string
}

func (u *User) Notify() error {
	fmt.Printf("User: Sending User Email To %s<%s>\n",
		u.Name,
		u.Email)
	return nil
}

func (u *User) Print() string {
	return fmt.Sprintf("Name = %s, Email = %s", u.Name, u.Email)
}

type Admin struct {
	Notifier
	Level string
}

func (u *Admin) Print() string {
	return u.Notifier.Print() + fmt.Sprintf(" Level = %s", u.Level)
}

func main() {
	user := &User{
		Name:  "janet jones",
		Email: "janet@email.com",
	}
	admin := &Admin{
		Notifier: user,
		Level:    "super",
	}
	user.Notify()
	fmt.Println(user.Print())

	admin.Notify()
	fmt.Println(admin.Print())
}
