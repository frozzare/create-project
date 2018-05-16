package project

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/otiai10/copy"
	"github.com/pkg/errors"
	"gopkg.in/AlecAivazis/survey.v1"
	"gopkg.in/src-d/go-git.v4"
)

const projectFilename = "project.json"

// File represents a project file.
type File struct {
	Lables map[string]string      `json:"labels"`
	Fields map[string]interface{} `json:"fields"`
	Meta   struct {
		Name string
	} `json:"meta"`
}

// Empty returns true if a file struct is empty.
func (f *File) Empty() bool {
	return len(f.Fields) == 0 && len(f.Lables) == 0
}

// Project represents a project.
type Project struct {
	file *File
	log  *log.Logger
	dst  string
	src  string
}

// New create a new project with the given options.
func New(options ...Option) *Project {
	var v Project
	for _, o := range options {
		o(&v)
	}

	if v.log == nil {
		v.log = log.New(os.Stdout, "", 0)
	}

	return &v
}

// Create creates a new project.
func (p *Project) Create() error {
	if err := p.clone(); err != nil {
		return err
	}

	if _, err := os.Stat(filepath.Join(p.dst, projectFilename)); os.IsNotExist(err) {
		return nil
	}

	if err := p.readFile(); err != nil {
		return err
	}

	if err := filepath.Walk(p.dst, p.process); err != nil {
		return err
	}

	return p.clean()
}

func (p *Project) readFile() error {
	buf, err := ioutil.ReadFile(filepath.Join(p.dst, projectFilename))
	if err != nil {
		return err
	}

	err = json.Unmarshal(buf, &p.file)

	if err != nil || p.file == nil || p.file.Empty() {
		var fields map[string]interface{}
		if err := json.Unmarshal(buf, &fields); err != nil {
			return err
		}
		p.file = &File{Fields: fields}
	}

	var mk []string
	for k := range p.file.Fields {
		mk = append(mk, k)
	}
	sort.Strings(mk)

	for _, k := range mk {
		s, err := p.read(k, p.file.Fields[k])
		if err != nil {
			return err
		}

		p.file.Fields[k] = s
	}

	return nil
}

func (p *Project) read(key string, value interface{}) (string, error) {
	switch v := value.(type) {
	case []interface{}:
		return p.readSelect(key, v)
	default:
		s := strings.TrimSpace(fmt.Sprintf("%v", v))

		s, err := renderTemplate(s, p.file.Fields)
		if err != nil {
			return "", err
		}

		return p.readString(key, s)
	}
}

func (p *Project) readString(key, value string) (string, error) {
	var output string

	if p.file.Lables != nil && p.file.Lables[key] != "" {
		key = p.file.Lables[key]
	}

	err := survey.AskOne(&survey.Input{
		Message: key,
		Default: value,
	}, &output, nil)

	if err.Error() == "EOF" {
		output = value
		err = nil
	}

	return strings.TrimSpace(output), err
}

func (p *Project) readSelect(key string, value []interface{}) (string, error) {
	var output string
	var options []string

	for _, v := range value {
		options = append(options, fmt.Sprintf("%v", v))
	}

	if p.file.Lables != nil && p.file.Lables[key] != "" {
		key = p.file.Lables[key]
	}

	err := survey.AskOne(&survey.Select{
		Message: key,
		Options: options,
	}, &output, nil)

	if err.Error() == "EOF" {
		if len(options) == 0 {
			output = ""
		} else {
			output = options[0]
		}

		err = nil
	}

	return strings.TrimSpace(output), err
}

func (p *Project) clone() error {
	if strings.Contains(p.src, "http") {
		_, err := git.PlainClone(p.dst, false, &git.CloneOptions{
			URL: p.src,
		})
		return err
	}

	return copy.Copy(p.src, p.dst)
}

func (p *Project) process(oldPath string, fi os.FileInfo, err error) error {
	if err != nil {
		if !strings.Contains(strings.ToLower(err.Error()), "lstat") {
			return err
		}
	}

	if oldPath[0] == '.' {
		return nil
	}

	path, err := renderTemplate(oldPath, p.file.Fields)
	if err != nil {
		return errors.Wrap(err, "path template")
	}

	if path == oldPath {
		return nil
	}

	if _, err := os.Stat(oldPath); err == nil {
		if err := os.Rename(oldPath, path); err != nil {
			return err
		}
	}

	if fi == nil {
		file, err := os.Open(path)
		if err != nil {
			return err
		}

		fi, err = file.Stat()
		if err != nil {
			return err
		}
	}

	if fi.IsDir() {
		return nil
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return errors.Wrap(err, "read file")
	}

	s, err := renderTemplate(string(buf), p.file.Fields)
	if err != nil {
		return errors.Wrap(err, "file template")
	}

	return ioutil.WriteFile(path, []byte(s), fi.Mode())
}

func (p *Project) clean() error {
	return os.Remove(filepath.Join(p.dst, projectFilename))
}
