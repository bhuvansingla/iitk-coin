package mail

import (
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

func Test() (err error) {
	authorize()
	to := []string{from}
	msg := []byte("To: " + from + "\n" +
		"From: " + "IITK-Coin<" + from + ">\n" +
		"Subject: IITK-Coin Test Mail\n" +
		"This is a Test Mail" + 
		"\n")
		
	err = smtp.SendMail(server.Address(), auth, from, to, msg)
	if err != nil {
		return
	}
	log.Info("Test mail sent")
	return
}
