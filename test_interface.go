package main

import (
	"fmt"
	"unsafe"
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

type eface struct {
	typ, val unsafe.Pointer
}

type poolDequeue struct {
	head, tail uint32
	vals       []eface
}

func (d *poolDequeue) push(val interface{}) {
	head := d.head
	fmt.Println("queue head: ", d.head)
	slot := &d.vals[head&uint32(len(d.vals)-1)]
	*(*interface{})(unsafe.Pointer(slot)) = val
	d.head = (head + 1) % (uint32(len(d.vals)) - 1)
}

func interfaceIsNil(x interface{}) {
	if x == nil {
		fmt.Println("empty interface")
		return
	}
	fmt.Println("non-empty interface")
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

	fmt.Printf("\n\n")
	const initSize = 8
	queue := new(poolDequeue)
	queue.vals = make([]eface, initSize)

	queue.push(user)
	queue.push(1)
	queue.push(23.1)
	fmt.Printf("%#v\n", queue)
	fmt.Printf("%#v, %#v\n\n", queue.vals[0].typ, queue.vals[0].val)

	var x interface{} = nil
	var y *int = nil
	interfaceIsNil(x)
	interfaceIsNil(y)
}
