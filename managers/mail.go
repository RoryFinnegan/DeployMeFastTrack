package manager

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
	"time"
)

// Variable to store information about send/receiving mail
var Mail MailConfig

// Generate a unique message ID for the sent email in
// accordance with RFC 5322
func generateMessageID() string {
	var id [12]byte
	_, err := rand.Read(id[:])
	if err != nil {
		return fmt.Sprintf("%d.%s@%s", time.Now().UnixNano(), "default", Mail.Server)
	}
	return fmt.Sprintf("%x@%s", id[:], Mail.Server)
}

// Retrieve all information regaurding email from
// config.json
func MailGetConfig() error {
	file, err := os.Open("config.json")
	if err != nil {
		return err
	}
	defer file.Close()

	//Decode the json file into a GoLang struct
	decoder := json.NewDecoder(file)
	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		return err
	}

	Mail = config.Mail

	return nil
}

// Uses TLS encryption to pass user login information
// to the mail server, and pass a message indicating
// an addition to the asset list
func SendMail(asset string, snum string, user string) {
	auth := smtp.PlainAuth("", Mail.Sender, Mail.Password, Mail.Server)
	messageID := generateMessageID()
	to := []string{Mail.Receiver}
	msg := []byte(
		"To: " + Mail.Receiver + "\r\n" +
			"From: " + Mail.Sender + "\r\n" +
			"Subject: [" + asset + "] Deployed by " + user + "\r\n" +
			"Message-ID: <" + messageID + ">\r\n" +
			"\r\n" +
			"The asset [" + asset + "] with serial [" + snum + "] has been taken by " + user + "\r\n")
	err := smtp.SendMail(Mail.Server+":"+Mail.Port, auth, Mail.Sender, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
