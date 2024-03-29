package events

import "log"

type Event interface {
	GetName() string
	Run()
}

type Function struct {
	Name     string
	Function func()
}

func (f Function) GetName() string {
	return f.Name
}
func (f Function) Run() {
	f.Function()
}

func FuncEvent(name string, function1 func()) {

	log.Printf(" <func %s> %s", name, "Begin")
	function1()
	log.Printf(" <func %s> %s", name, "Finish")
}
func GoFuncEvent(name string, function1 func()) {
	go func() {
		log.Printf(" <go %s> %s", name, "Begin")
		function1()
		log.Printf(" <func %s> %s", name, "Finish")
	}()
}
func HandleEvent(event Event) {
	log.Printf(" <%s> %s", event.GetName(), "Begin")
	event.Run()
	log.Printf(" <%s> %s", event.GetName(), "Finish")
}

func DoneFuncEvent(name string, function1 func(chan bool), Shutdown chan bool) {
	log.Printf(" <func(%s)> %s\n", name, " Starting Done Func")
	function1(Shutdown)
}
