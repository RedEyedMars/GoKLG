package networking

import (
	"strings"

	"../events"
)

var adminCommands map[string]events.Event
var adminArgs []string

func SetupAdminCommands() {
	if adminCommands == nil {
		adminCommands = make(map[string]events.Event)
		adminCommands["exit"] = &events.Function{Name: "Admin!Exit", Function: func() { Shutdown <- true }}
		/*adminCommands["addMember"] = &events.Function{Name: "Admin!AddMember", Function: func() {
			if adminArgs != nil {
				memberIp := adminArgs[0]
				databasing.NewMember(memberIp)
			}
		}}
		*/
	}
}
func HandleAdminCommand(msg string) bool {

	splice := strings.Split(msg, " ")
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
