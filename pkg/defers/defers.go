// Package defers allows to have "global defer", by pushing functions into this
// packages data holder, which can then be ran from main
package defers

var defers []func()

// Run should be deferred from the most outer function, to allow "global defer"
func Run() {
	// last in. first out
	for ii := range defers {
		i := len(defers) - 1 - ii
		defers[i]()
	}
}

// Defer defers a function globally. These typically run at program exit
func Defer(f func()) {
	defers = append(defers, f)
}
