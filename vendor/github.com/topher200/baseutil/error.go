// Utility functions for all golang needs

package baseutil

// Panic if error is non-nil.
//
// Useful for functions that return an error that you never expect them to
// throw. Just Check(err) after their call and move on.
func Check(e error) {
	if e != nil {
		panic(e)
	}
}
