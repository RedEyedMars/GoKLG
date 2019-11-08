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
			c.id = user.ID
			c.send <- []byte(fmt.Sprintf("{login_successful:%s}", user.Name))
			log.Printf(" networking.attempt_login.Login successful")
		} else {
			c.send <- []byte("{login_failed}Credentials not accepted, either check your password or your username!")
			log.Printf(" networking.attempt_login.Login failed")
		}
	}
	commands["attempt_signup"] = func(c *Client, msg []byte, user []byte) {
		username := string(user)
		if user := <-databasing.RequestUser("ByName", username); user != nil {
			c.send <- []byte("{signup_failed}Username taken!")
		} else {

			log.Printf(" networking.attempt_signup.No user found; good!")
			hash := sha256.New()
			hash.Write([]byte(adminPassword))
			hash.Write(msg) //Add Password here
			pwdAsString := fmt.Sprintf("%x", hash.Sum(nil)[:])
			if user := <-databasing.RequestUser("ByPwd", pwdAsString); user == nil {

				events.FuncEvent("client.Signup.AddUser", func() {
					user = <-databasing.InsertUser(username, pwdAsString)
				})
				c.name = user.Name
				c.id = user.ID

				c.send <- []byte(fmt.Sprintf("{signup_successful:%s}", user.Name))
				log.Printf(" networking.attempt_signup.Signup Successful!")
			} else {
				c.send <- []byte("{login_failed}Credentials not accepted, try a different password and username!")
				log.Printf(" networking.attempt_signup.Signup Failed!")
			}
		}
	}

	commands["attempt_logout"] = func(c *Client, msg []byte, user []byte) {
		c.name = "!none!"
		c.id = -1
		c.send <- []byte("{logout_successful}")
	}

	commands["attempt_change"] = func(c *Client, msg []byte, user []byte) {
		msgSplit := strings.Split(string(msg), ",")
		newUser, oldPwd, pwd := msgSplit[0], msgSplit[1], msgSplit[2]
		hash := sha256.New()
		hash.Write([]byte(adminPassword))
		hash.Write([]byte(pwd))
		pwdAsString := fmt.Sprintf("%x", hash.Sum(nil)[:])

		hash = sha256.New()
		hash.Write([]byte(adminPassword))
		hash.Write([]byte(oldPwd))
		oldPwdAsString := fmt.Sprintf("%x", hash.Sum(nil)[:])
		if confirmed_user := <-databasing.RequestUser("ByPwd", oldPwdAsString); confirmed_user != nil {
			if success := <-databasing.ChangeUser(string(user), newUser, pwdAsString); success {
				c.name = newUser
				c.send <- []byte(fmt.Sprintf("{change_successful:%s}", newUser))
				log.Printf(" networking.attempt_change.Change successful")
			} else {
				c.send <- []byte("{change_failed}No matching username was found for:" + string(user) + "!")
				log.Printf(" networking.attempt_change.Change failed")
			}
		} else {
			c.send <- []byte("{change_failed}No matching username was found for:" + string(user) + "!")
			log.Printf(" networking.attempt_change.Change failed")
		}
	}
}
