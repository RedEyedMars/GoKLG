package Events

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

func FuncEvent(name string, Function1 func()) {
	HandleEvent(Function{Name: "func(" + name + ")", Function: Function1})
}
func GoFuncEvent(name string, Function1 func()) {
	go HandleEvent(Function{Name: "go func(" + name + ")", Function: Function1})
}
func HandleEvent(event Event) {
	log.Printf(" <%s> %s", event.GetName(), "Begin")
	event.Run()
	log.Printf(" <%s> %s", event.GetName(), "Finish")
}

func DoneFuncEvent(name string, Function1 func(chan bool), Shutdown chan bool) {
	log.Printf(" <%s> %s", "func("+name+")", "Finish")
	Function1(Shutdown)
}
