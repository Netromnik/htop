package top

var RegisterProgram *ProgramPool

var Observer *MultiChanObserver

func init() {
	Observer = NewMultiChanObserver()
	RegisterProgram = newProgramPool()
	newFileHealchek("active", RegisterProgram)
}

// func main() {
// 	Observer = NewMultiChanObserver()
// 	RegisterProgram = newProgramPool()
// 	fh := newFileHealchek("active", RegisterProgram)
// 	fh.load()
// 	go fh.wath()

// 	ch := Observer.AddChannel("test")
// 	fmt.Println(len(*ch), <-*ch)
// }
