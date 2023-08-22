package top

import (
	"fmt"
	"testing"
)

func TestLoadFileWath(t *testing.T) {
	Observer = NewMultiChanObserver()
	poll := newProgramPool()
	// fh := newFileHealchek("active", poll)
	// fh.run()
	fmt.Println(poll.list_app)
}
