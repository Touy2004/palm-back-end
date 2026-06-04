package main

import (
	"fmt"
	"go/parser"
	"go/token"
)

func main() {
	fset := token.NewFileSet()
	_, err := parser.ParseFile(fset, "internal/middleware/role.go", nil, parser.AllErrors)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Parsed successfully!")
	}
}
