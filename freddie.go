package main

import ()

type Freddie struct {
}

func newFreddie() Freddie {
	return Freddie{}
}
func (f Freddie) expose() string {
	return "# HELP freddie A simple header for this exporter.\n# TYPE freddie counter\nfreddie 1\n"
}
