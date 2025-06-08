package main

import "fmt"

func main() {
	project, err := Parse("./test-extracted/project.json")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v\n", project)
}
