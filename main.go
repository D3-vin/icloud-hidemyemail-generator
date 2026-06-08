package main

import (
	"os"
	"github.com/yourusername/hidemyemail-generator/cmd/hidemyemail"
)

func main() {
	hidemyemail.Execute()
	os.Exit(0)
}
