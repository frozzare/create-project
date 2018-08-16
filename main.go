package main

import (
	"log"
	"os"
	"strings"

	"github.com/frozzare/create-project/project"
)

var version = "master"

func main() {
	log.SetFlags(0)
	log.SetOutput(os.Stderr)

	if len(os.Args) < 2 {
		log.Fatal("No source url or directory")
	}

	src := os.Args[1]

	if strings.ToLower(src) == "version" {
		log.Printf("create-project version %s\n", version)
		return
	}

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
