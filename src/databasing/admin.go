package databasing

import (
	"fmt"
	"log"
	"strings"

	"../events"
)

var adminCommands map[string]events.Event
var adminArgs []string

func MakeAdminFunc(argCount uint16, f func(...string)) func() {
	switch argCount {
	case 0:
		return func() { f() }
	case 1:
		return func() {
			if adminArgs != nil && len(adminArgs) > 0 {
				f(adminArgs[0][:len(adminArgs[0])-1])
			}
		}
	case 2:
		return func() {
			if adminArgs != nil && len(adminArgs) > 1 {
				f(adminArgs[0], adminArgs[1][:len(adminArgs[1])-1])
			}
		}
	case 3:
		return func() {
			if adminArgs != nil && len(adminArgs) > 2 {
				f(adminArgs[0], adminArgs[1], adminArgs[2][:len(adminArgs[2])-1])
			}
		}
	}
	return func() {}
}
func SetupAdminCommands() {
	if adminCommands == nil {
		adminCommands = make(map[string]events.Event)
		//adminCommands["exit"] = &events.Function{Name: "Admin!Exit", Function: func() { Shutdown <- true }}
		adminCommands["add"] = &events.Function{Name: "Admin!Add", Function: MakeAdminFunc(2,
			func(args ...string) { InsertUser(args[0], args[1]) })}
		adminCommands["remove"] = &events.Function{Name: "Admin!Remove", Function: MakeAdminFunc(1,
			func(args ...string) { DeleteUser(args[0]) })}
		adminCommands["change"] = &events.Function{Name: "Admin!Change", Function: MakeAdminFunc(3,
			func(args ...string) { ChangeUser(args[0], args[1], args[2]) })}

		adminCommands["report_oct"] = &events.Function{Name: "Admin!Report", Function: MakeAdminFunc(3,
			func(args ...string) {
				result := <-RequestOctoberReport()
				fmt.Println(result)
			})}
	}
}
func HandleAdminCommand(msg string) bool {
	splice := strings.Split(msg, " ")
	log.Printf(" database handle :" + msg)
	if len(splice) == 1 {
		if command := adminCommands[msg]; command == nil {
			return false
		} else {
			events.HandleEvent(command)
			return true
		}
	} else {
		if command := adminCommands[splice[0]]; command == nil {
			return false
		} else {
			adminArgs = splice[1:]
			events.HandleEvent(command)
			return true
		}
	}
}
