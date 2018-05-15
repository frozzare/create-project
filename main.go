package main

import (
	"log"
	"os"

	"github.com/frozzare/create-project/project"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("No source url or directory")
	}

	src := os.Args[1]

	path, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dst := path
	if len(os.Args) > 2 {
		dst = os.Args[2]
	}

	p := project.New(
		project.Source(src),
		project.Destination(dst),
	)

	if err := p.Create(); err != nil {
		log.Fatal(err)
	}
}
