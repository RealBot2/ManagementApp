package main

import (
	"fmt"
)

func main() {
	db := CheckInitDB()
	defer db.Close()

	fmt.Println("Hello")
}
