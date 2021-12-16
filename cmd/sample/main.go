package main

import "fmt"

type Keeper struct {
	n Nested
}

type Nested struct {
	i int
	m map[string]int
}

func (_k Keeper) SetMe(k string, v int) {
	_k.n.m[k] = v
}

func (_k Keeper) Get(k string) int {
	return _k.n.m[k]
}

func (_k Keeper) SetI(v int) {
	_k.n.i = v
}

func (_k Keeper) GetI() int {
	return _k.n.i
}

func main() {
	k := Keeper{
		n: Nested{
			i: 0,
			m: make(map[string]int),
		},
	}

	k.SetMe("a", 1)

	fmt.Printf("%+v \n", k.Get("a"))

	k.SetI(1)

	fmt.Printf("%+v \n", k.GetI())

}
