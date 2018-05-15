package project

import "log"

// Option configures a project option.
type Option func(*Project)

// Destination returns an option configured with a destination value.
func Destination(v string) Option {
	return func(p *Project) {
		p.dst = v
	}
}

// Source returns an option configured with a source value.
func Source(v string) Option {
	return func(p *Project) {
		p.src = v
	}
}

// Log returns an option configured with a log value.
func Log(v *log.Logger) Option {
	return func(p *Project) {
		p.log = v
	}
}
