package mail

import (
	"net/smtp"

	log "github.com/sirupsen/logrus"
)

func mailService(mailChannel chan EmailRequest) {
	authorize()
	for request := range mailChannel {
		to := []string{request.To+"@iitk.ac.in"}
		msg := []byte("To: " + request.To + "@iitk.ac.in" + "\n" +
			"From: " + "IITK-Coin<" + from + ">\n" +
			"Subject: IITK-Coin One Time Password\n" +
			"Your OTP is " + request.OTP + "\n" +
			"This OTP is valid for " + otpValidity + " minutes and don't share it with anyone." + 
			"\n")
			
		err := smtp.SendMail(server.Address(), auth, from, to, msg)

		// if error, try to login again
		if err != nil {
			authorize()
			err = smtp.SendMail(server.Address(), auth, from, to, msg)
			if err != nil {
				log.Error("Error sending mail: " + err.Error())
				continue
			}
		}
		
		log.Info("Mail sent to ", request.To)
	}
}
