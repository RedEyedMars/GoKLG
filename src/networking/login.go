package networking

import (
	"crypto/sha256"
	"fmt"
	"log"
	"strings"

	"../databasing"
	"../events"
)

func setupLoginCommands(registry *ClientRegistry) {

	commands["attempt_login"] = func(c *Client, msg []byte, user []byte) {
		hash := sha256.New()
		hash.Write([]byte(adminPassword))
		hash.Write(msg)
		pwdAsString := fmt.Sprintf("%x", hash.Sum(nil)[:])
		if user := <-databasing.RequestUser("ByPwd", pwdAsString); user != nil {
			c.name = user.Name
			c.send <- []byte(fmt.Sprintf("{login_successful:%s}", user.Name))
			log.Printf(" networking.attempt_login.Login successful")
		} else {
			c.send <- []byte("{login_failed}Credentials not accepted, either check your password or your username!")
			log.Printf(" networking.attempt_login.Login failed")
		}
	}
	commands["attempt_signup"] = func(c *Client, msg []byte, user []byte) {
		split := strings.Split(string(msg), ",")
		username, pwd := split[0], split[1]
		if user := <-databasing.RequestUser("ByName", username); user != nil {
			c.send <- []byte("{signup_failed}Username taken!")
		} else {

			log.Printf(" networking.attempt_signup.No user found; good!")
			hash := sha256.New()
			hash.Write([]byte(adminPassword))
			hash.Write([]byte(pwd))
			pwdAsString := fmt.Sprintf("%x", hash.Sum(nil)[:])
			if user := <-databasing.RequestUser("ByPwd", pwdAsString); user == nil {

				events.FuncEvent("client.Signup.AddUser", func() {
					user = <-databasing.InsertUser(username, pwdAsString)
				})
				c.name = user.Name

				c.send <- []byte(fmt.Sprintf("{signup_successful:%s}", user.Name))
				log.Printf(" networking.attempt_signup.Signup Successful!")
			} else {
				c.send <- []byte("{login_failed}Credentials not accepted, try a different password and username!")
				log.Printf(" networking.attempt_signup.Signup Failed!")
			}
		}
	}

	commands["attempt_logout"] = func(c *Client, msg []byte, user []byte) {
		c.name = "_none_"
		c.send <- []byte("{logout_successful}")
	}
}
