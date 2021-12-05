package mail

import (
	"net/smtp"

	"github.com/spf13/viper"
)

type smtpServer struct {
	host string
	port string
}

func (s *smtpServer) Address() string {
	return s.host + ":" + s.port
}

type EmailRequest struct {
	To  string
	OTP string
}

var (
	MailChannel chan EmailRequest
	from        string
	password    string
	server      smtpServer
	auth        smtp.Auth
	otpValidity string
)

func init() {
	MailChannel = make(chan EmailRequest)
	from = viper.GetString("MAIL.FROM")
	password = viper.GetString("MAIL.PASSWORD")
	
	otpValidity = viper.GetString("OTP.EXPIRY_PERIOD_IN_MIN")

	server = smtpServer{host: viper.GetString("MAIL.HOST"), port: viper.GetString("MAIL.PORT")}
	go mailService(MailChannel)
}

func authorize() {
	auth = smtp.PlainAuth("", from, password, server.host)
}
