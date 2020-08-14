package main

import (
	"os"
	"test/app"
)

func main() {
	os.Exit(app.New().Run(os.Args[1:]))
}
