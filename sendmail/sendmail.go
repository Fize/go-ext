package sendmail

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"

	"github.com/fize/go-ext/log"
)

func dial(addr string) (*smtp.Client, error) {
	conn, err := tls.Dial("tcp", addr, nil)
	if err != nil {
		log.Warnf("connect to smtp server %s error", addr)
		return nil, err
	}
	// parse host and port string
	host, _, _ := net.SplitHostPort(addr)
	return smtp.NewClient(conn, host)
}

func sendMailUsingTLS(addr string, auth smtp.Auth, from string,
	to []string, msg []byte) (err error) {

	// create smtp client
	c, err := dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()

	if auth != nil {
		if ok, _ := c.Extension("AUTH"); ok {
			if err = c.Auth(auth); err != nil {
				return err
			}
		}
	}

	if err = c.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	_, err = w.Write(msg)
	if err != nil {
		return err
	}

	err = w.Close()
	if err != nil {
		return err
	}

	return c.Quit()
}

// SendEmail send email
func SendEmail(headerTitle, smtpServer, email, password, to, subject, body string, port int) error {
	header := make(map[string]string)
	header["From"] = headerTitle + "<" + email + ">"
	header["To"] = to
	header["Subject"] = subject
	header["Content-Type"] = "text/html; charset=UTF-8"

	message := ""
	for k, v := range header {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + body

	auth := smtp.PlainAuth(
		"",
		email,
		password,
		smtpServer,
	)

	err := sendMailUsingTLS(
		fmt.Sprintf("%s:%d", smtpServer, port),
		auth,
		email,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return err
	}
	log.Infof("send email to %s, content: %s", to, body)
	return nil
}
