package databasing

import (
	"log"
	"strings"

	"../events"
)

var adminCommands map[string]events.Event
var adminArgs []string

func makeAdminFunc(argCount uint16, f func(...string)) func() {
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
	}
	return func() {}
}
func SetupAdminCommands() {
	if adminCommands == nil {
		adminCommands = make(map[string]events.Event)
		//adminCommands["exit"] = &events.Function{Name: "Admin!Exit", Function: func() { Shutdown <- true }}
		adminCommands["addMember"] = &events.Function{Name: "Admin!AddMember_Full", Function: makeAdminFunc(2,
			func(args ...string) { InsertUser(args[0], args[1]) })}
		adminCommands["removeMember"] = &events.Function{Name: "Admin!RemoveMember", Function: makeAdminFunc(1,
			func(args ...string) { RequestAction("Users", "Remove", args[0]) })}
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
