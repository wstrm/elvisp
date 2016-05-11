package tasks

const (
	errorInvalidLength = "2 Invalid length of arguments"
	errorInvalidTask   = "1 Invalid task specified"
)

// TaskInterface defines the methods needed for a default task
type TaskInterface interface {
	Run() string
}

// Task contains the arguments that the task needs
type Task struct {
	argv []string
}

// SetArgs sets the argument array
func (t Task) SetArgs(a []string) {
	t.argv = a
}

func (t Task) errorString(e string) string {
	return "error " + e
}

func (t Task) successString(s string) string {
	return "success " + s
}

// Add should implement the add task
type Add struct{ Task }

// Remove should implement the remove task
type Remove struct{ Task }

// Lease should implement the lease task
type Lease struct{ Task }

// Release should implement the release task
type Release struct{ Task }

// Invalid should implement an invalid task which is returned if no task could be found for a command
type Invalid struct{ Task }

// Run adds a user using the public key and a token
func (t Add) Run() string {
	if len(t.argv) < 3 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("<<LEASED ADDRESS HERE>>")
}

// Run ...
func (t Remove) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("Removed user: " + t.argv[0])
}

// Run ...
func (t Lease) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("<<LEASED ADDRESS HERE>>")
}

// Run ...
func (t Release) Run() string {
	if len(t.argv) != 2 {
		return t.errorString(errorInvalidLength)
	}

	return t.successString("Released lease for user: " + t.argv[0])
}

// Run ...
func (t Invalid) Run() string {
	return t.errorString(errorInvalidTask)
}
