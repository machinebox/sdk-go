package boxutil

// Box represents a box client capable of returning
// Info.
type Box interface {
	Info() (*Info, error)
}

// Info describes box information.
type Info struct {
	Name    string
	Version int
	Build   string
	Status  string
}
