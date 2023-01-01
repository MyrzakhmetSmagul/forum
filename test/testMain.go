package main

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	start := time.Now()
	password := os.Args[1]
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(password)
	fmt.Println(string(hash))
	// fmt.Println(bcrypt.CompareHashAndPassword(hash, []byte(password)))
	fmt.Println(time.Since(start))
}
