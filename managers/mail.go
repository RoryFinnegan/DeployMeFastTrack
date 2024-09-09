package manager

import (
	"encoding/json"
	"fmt"
	"net/smtp"
	"os"
)

var Mail MailConfig

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

func SendMail(asset string, snum string, user string) {
	auth := smtp.PlainAuth("", Mail.Sender, Mail.Password, Mail.Server)
	to := []string{Mail.Receiver}
	msg := []byte("To: " + Mail.Receiver + "\r\n" +
		"Subject: [" + asset + "] Deployed by " + user + "\r\n" +
		"\r\n" +
		"The asset [" + asset + "] with serial [" + snum + "] has been taken by " + user + "\r\n")
	err := smtp.SendMail(Mail.Server+":"+Mail.Port, auth, Mail.Sender, to, msg)
	if err != nil {
		fmt.Println(err)
		return
	}
}
