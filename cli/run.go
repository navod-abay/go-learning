package main

import (
	generate "github.com/navod-abay/mandelbrotset-go/core"
)

// Dummy module to run the cli tool by itself while importing from the core module
func main() {
	generate.Run()
}
