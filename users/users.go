package users

import (
	"errors"
	"os"
	"text/template"

	"github.com/quanvuminh/gastwrap/ari"

	scribble "github.com/nanobox-io/golang-scribble"
)

const (
	userConfTempl  string = "templates/userid.tmpl"
	userConfPath   string = "/etc/asterisk/endpoints/"
	defaultContext string = "hapulico-internal"
	dbPath         string = "/etc/asterisk/endpoints/userjson"
)

// User define an endpoint
type User struct {
	UserID     string `json:"userid" form:"userid"`
	CallerID   string `json:"callerid" form:"callerid"`
	Password   string `json:"password" form:"password"`
	Context    string `json:"context" form:"context"`
	PhoneModel string `json:"phonemodel" form:"phonemodel"`
	PhoneMac   string `json:"phonemac" form:"phonemac"`
}

// Find an existing user
func Find(userid string) (status string) {
	userconf := userConfPath + userid + ".conf"
	if _, err := os.Stat(userconf); os.IsNotExist(err) {
		return "Available"
	}

	return "Unavailable"
}

func writeToDB(user User) error {
	db, errNew := scribble.New(dbPath, nil)
	if errNew != nil {
		return errNew
	}

	err := db.Write("user", user.UserID, user)

	return err
}

// PjsipReload run alias "pjsip reload"
func PjsipReload() error {
	listModules := []string{
		"res_pjsip.so",
		"res_pjsip_authenticator_digest.so",
		"res_pjsip_endpoint_identifier_ip.so",
		"res_pjsip_mwi.so res_pjsip_notify.so",
		"res_pjsip_outbound_publish.so",
		"res_pjsip_publish_asterisk.so",
		"res_pjsip_outbound_registration.so",
		"res_phoneprov.so",
	}
	for _, m := range listModules {
		if err := ari.ModuleReload(m); err != nil {
			return err
		}
	}

	return nil
}

// Create a new user
func Create(newuser User) error {
	// Check if the new user exists
	if stt := Find(newuser.UserID); stt == "Unavailable" {
		return errors.New("Exists")
	}

	if newuser.CallerID == "" {
		newuser.CallerID = newuser.UserID
	}

	if newuser.Context == "" {
		newuser.Context = defaultContext
	}

	tmpl, err := template.ParseFiles(userConfTempl)
	if err != nil {
		return err
	}

	userconf := userConfPath + newuser.UserID + ".conf"
	f, err := os.OpenFile(userconf, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	err = tmpl.Execute(f, newuser)
	if err != nil {
		return err
	}

	if err = writeToDB(newuser); err != nil {
		return err
	}

	// Check if the new user is created
	if stt := Find(newuser.UserID); stt == "Unavailable" {
		return nil
	}

	if err = PjsipReload(); err != nil {
		return err
	}

	return errors.New("Failed")
}

// Get user infos
func Get(userid string) (user User, err error) {
	db, _ := scribble.New(dbPath, nil)
	err = db.Read("user", userid, &user)

	return user, err
}

// Delete an user
func Delete(userid string) error {
	// Delete user in DB
	db, _ := scribble.New(dbPath, nil)
	if err := db.Delete("user", userid); err != nil {
		return err
	}

	// Delete user configuration file
	userconf := userConfPath + userid + ".conf"
	if err := os.Remove(userconf); err != nil {
		return err
	}

	err := PjsipReload()

	return err
}

// Update user infos
func Update(olduserid string, newuser User) error {
	// Get current values
	u, e := Get(olduserid)
	if e != nil {
		return e
	}

	if newuser.UserID == "" {
		newuser.UserID = u.UserID
	}
	if newuser.CallerID == "" {
		newuser.CallerID = u.CallerID
	}
	if newuser.Password == "" {
		newuser.Password = u.Password
	}
	if newuser.Context == "" {
		newuser.Context = u.Context
	}
	if newuser.PhoneModel == "" {
		newuser.PhoneModel = u.PhoneModel
	}
	if newuser.PhoneMac == "" {
		newuser.PhoneMac = u.PhoneMac
	}

	// Delete old infos
	if err := Delete(olduserid); err != nil {
		return err
	}

	if err := Create(newuser); err != nil {
		return err
	}

	err := PjsipReload()

	return err
}
