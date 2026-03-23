package main

import "fmt"

type Person struct {
	name string
	age  int
}

func changeName(p *Person, name string) {
	p.name = name
}

func getPerson() Person {
	var p Person

	defer changeName(&p, "Alice")

	p = Person{
		name: "Bob",
		age:  52,
	}

	return p
}

func getPerson2() (p Person) {
	defer changeName(&p, "Alice")

	p = Person{
		name: "Bob",
		age:  52,
	}

	return p
}

func main() {
	fmt.Println("main start")
	p := getPerson()
	p2 := getPerson2()

	fmt.Println("main end", p, p2)
}
