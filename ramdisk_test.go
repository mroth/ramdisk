package ramdisk_test

import "testing"

// DONE test to verify the implementation for current platform satisfies the interface
// dont need to anymore, since it will barf due to typechecking now

/* basic test sequence I'm thinking about
- try to create, no errors
- use go FS to see if path exists
- if possible, stat the directory and try to read capacity?
- use go FS to write a file to ramdisk
- use go FS to read a file back from ramdisk, verify contents
- unmount ramdisk, no errors
- go FS to see if file exists (should not)
- go FS to verify mountpath no longer exists
*/
func TestCreate(t *testing.T) {
	t.Skip("pending") // TODO pending test
}

func TestDestroy(t *testing.T) {
	t.Skip("pending") // TODO pending test
}
