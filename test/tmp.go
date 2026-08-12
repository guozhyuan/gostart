package test

import (
	"fmt"
)

type Human interface {
	ChangeName(s string)
}
type Person struct {
	Name string
}

func (p *Person) ChangeName(s string) {
	p.Name = s
}

// func (p Person) ChangeName(s string) {
// 	p.Name = s
// }

func Func() {
	p1 := Person{"初始名"}
	p1.ChangeName("P1")
	fmt.Println(p1.Name)

	p2 := &Person{"初始名"}
	p2.ChangeName("P2")
	fmt.Println(p2.Name)
}

// pointer, channel, func, interface, map, or slice   == nil
func Slice() {

	// var slice []Person
	// fmt.Println(slice == nil)
	// fmt.Println(slice[0])
	// fmt.Println("--------------------------")
	var slice2 = make([]Person, 3)
	fmt.Println(slice2 == nil)
	fmt.Println(slice2[0])
	// fmt.Println("--------------------------")
	// var slice3 = []Person{}
	// fmt.Println(slice3 == nil)
	// fmt.Println(slice3[0])
}
