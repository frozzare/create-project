# create-project [![Build Status](https://travis-ci.org/frozzare/create-project.svg?branch=master)](https://travis-ci.org/frozzare/create-project) [![GoDoc](https://godoc.org/github.com/frozzare/create-project?status.svg)](http://godoc.org/github.com/frozzare/create-project) [![Go Report Card](https://goreportcard.com/badge/github.com/frozzare/create-project)](https://goreportcard.com/report/github.com/frozzare/create-project)

> Work in progress

Create project is a command line tool for create new project from a boilerplate.

## Installation

```
go get -u github.com/frozzare/create-project
```

Or download the release from release page.

## Usage

Create a `project.json` in your boilerplate directory.

```js
{
    "name": "simple"
}
```

Then create a directory called `{{.name}}` with a `main.js` file that contains this:

```js
var {{.name}} = function () {};
```

Then run:

```
create-project my-boilerplate dest-folder
```

You can also use a git url:

```
create-project https://github.com/user/my-boilerplate.git dest-folder
```

## Labels

You can add labels to all fields using a another JSON structure in your `project.json`

```js
{
    "fields": {
        "name": "app"
    },
    "labels": {
        "name": "Name"
    }
}
```

## Commands

Commands array can be added when using labels. The commands will be executed after the project is generated.

```js
{
    "fields": {
        "name": "app"
    },
    "labels": {
        "name": "Name"
    },
    "commands": [
        "touch test"
    ]
}
```

## License

MIT © [Fredrik Forsmo](https://github.com/frozzare)