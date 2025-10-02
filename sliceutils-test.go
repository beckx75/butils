package main

import (
	"fmt"

	"beckx.online/butils/sliceutils"
)

func main() {
	l1 := []string{"1", "2", "3", "2", "2", "4"}
	l2 := []string{"4", "2", "3", "4"}
	o1, o2, i, nu1, nu2 := sliceutils.CompareLists(l1, l2)
	fmt.Println("O1", o1)
	fmt.Println("O2", o2)
	fmt.Println("Inter", i)
	fmt.Println("NU1", nu1)
	fmt.Println("NU2", nu2)

}
