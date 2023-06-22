package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/dawkaka/go-interpreter/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s, this is the monkey programming language!\n", user.Username)
	fmt.Println("Feel free to key in your commands")
	repl.Start(os.Stdin, os.Stdout)
}
