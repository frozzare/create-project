package project

import (
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestProject(t *testing.T) {
	for _, src := range []string{"./testdata/project", "./testdata/project2"} {
		os.RemoveAll("./testdata/test")

		p := New(Source(src), Destination("./testdata/test"))

		if err := p.Create(); err != nil {
			t.Fatal(err)
		}

		buf, err := ioutil.ReadFile("./testdata/test/app/main.js")
		if err != nil {
			t.Fatal(err)
		}

		if strings.TrimSpace(string(buf)) != "var app = function () {};" {
			t.Fatal("Expected main.js to be equal")
		}
	}
}
