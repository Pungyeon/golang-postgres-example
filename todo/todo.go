package todo

// Todo struct descripting todo objects
type Todo struct {
	UID         int
	Title       string
	Description string
	Username    string // guid
	Completed   bool
}