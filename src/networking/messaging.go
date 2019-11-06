package networking

/*
/((([A-Za-z]{3,9}:(?:\/\/)?)(?:[-;:&=\+\$,\w]+@)?[A-Za-z0-9.-]+|(?:www.|[-;:&=\+\$,\w]+@)[A-Za-z0-9.-]+)((?:\/[\+~%\/.\w-_]*)?\??(?:[-\+=&;%@.\w_]*)#?(?:[\w]*))?)/
*/

import (
	"regexp"
	"strconv"
)

var reCurls *regexp.Regexp
var reAngles *regexp.Regexp
var reCommandMsg *regexp.Regexp
var reIPvPort *regexp.Regexp

func setupNetworkingRegex() {
	reCurls = regexp.MustCompile(`\\{([^\\}]+)\\}`)
	reAngles = regexp.MustCompile(`<([^>]+)>`)

	reCommandMsg = regexp.MustCompile(`\{([^\{\}:;]+)(:)?([a-zA-Z0-9_-]+)?\}(.*)`)
	reIPvPort = regexp.MustCompile(`([^:]+):(.+)`)
}
func DifferentiateMessage(incomingMsg []byte) (string, []byte, []byte) {
	result := reCommandMsg.FindSubmatch(incomingMsg)
	if len(result) == 5 {
		return string(result[1]), result[4], result[3]
	}
	return string(result[1]), result[2], nil
}
func getIPFromAddress(ipAddress string) (string, int) {

	result := reIPvPort.FindStringSubmatch(ipAddress)
	if port, err := strconv.Atoi(result[2]); err != nil {
		return result[1], -1
	} else {
		return result[1], port
	}
}
func SanatizeMessage(incomingMsg []byte) []byte {
	return reAngles.ReplaceAllFunc(reCurls.ReplaceAllFunc(incomingMsg,
		func(curl []byte) []byte {
			return []byte("{{" + string(curl) + "}}")
		}),
		func(angle []byte) []byte {
			return []byte("&lt" + string(angle) + "&gt")
		})
}

func ConstructMessage(header string, user []byte, msg []byte) []byte {
	if user != nil {
		return concatCopyPreAllocate([]byte("{"+header), []byte(":"), user, []byte("}"), msg)
	}
	return concatCopyPreAllocate([]byte("{"+header), []byte("}"), msg)

}

func concatCopyPreAllocate(slices ...[]byte) []byte {
	var totalLen int
	for _, s := range slices {
		totalLen += len(s)
	}
	tmp := make([]byte, totalLen)
	var i int
	for _, s := range slices {
		i += copy(tmp[i:], s)
	}
	return tmp
}
