package ami

import (
	"fmt"
	"log"

	ami "github.com/heltonmarx/goami/ami"
)

const (
	//AMHostPort is the remote address for the Ast manager
	AMHostPort string = "127.0.0.1:5038"
	// AMUser is the username for AMHostPort
	AMUser string = "mana"
	// AMSecret is the scret for AMUser
	AMSecret string = "manasecret"
)

// Callout with phone id
func Callout(number string, siptrunk string, context string) {
	socket, err1 := ami.NewSocket(AMHostPort)
	if err1 != nil {
		fmt.Printf("socket error: %v\n", err1)
		return
	}
	_, err2 := ami.Connect(socket)
	if err2 != nil {
		return
	}

	//Login
	uuid, _ := ami.GetUUID()
	err3 := ami.Login(socket, AMUser, AMSecret, "Off", uuid)
	if err3 != nil {
		log.Printf("login error (%v)\n", err3)
	}

	//Originate
	chandata := "Local/dial@" + context
	vardata := "callee=" + number + "trunk=" + siptrunk
	origdata := ami.OriginateData{
		Channel:  chandata,
		Context:  context,
		Exten:    "s",
		Priority: 1,
		Variable: vardata,
	}
	_, err4 := ami.Originate(socket, uuid, origdata)
	if err4 != nil {
		log.Printf("Originating error (%v)\n", err4)
	}

	// Logoff
	ami.Logoff(socket, uuid)
}
