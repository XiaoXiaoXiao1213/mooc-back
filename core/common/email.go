package common

import (
	"net/smtp"
	"strings"
)

func SendMail(to, subject, content string) error {
	user := "managementJamie@163.com"
	password := "QBPIVGKJGHLZAFPN"
	host := "smtp.163.com:25"
	auth := smtp.PlainAuth("", user, password, "smtp.163.com")
	body := `
    <html>
    <body>
    <h3>
    ` + content + `
    </h3>
    </body>
    </html>
    `
	msg := []byte("To: " + to + "\r\nFrom: " + user + "<" + user + ">\r\nSubject: " + subject + "\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n" + body)
	send_to := strings.Split(to, ";")
	err := smtp.SendMail(host, auth, user, send_to, msg)
	return err
}
