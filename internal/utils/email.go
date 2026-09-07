package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"

	"gitee.com/meinongyihe/travel-api/internal/config"

	"github.com/jordan-wright/email"
)

// Email 发送邮件给指定收件人（多个收件人用逗号分隔）
func Email(conf config.EmailConf, to, subject, body string) error {
	toList := strings.Split(to, ",")
	return send(conf, toList, subject, body)
}

// EmailDefault 使用配置中的默认收件人发送邮件
func EmailDefault(conf config.EmailConf, subject, body string) error {
	to := strings.Split(conf.To, ",")
	return send(conf, to, subject, body)
}

func send(conf config.EmailConf, to []string, subject, body string) error {
	auth := smtp.PlainAuth("", conf.From, conf.Secret, conf.Host)
	e := email.NewEmail()
	if conf.Nickname != "" {
		e.From = fmt.Sprintf("%s <%s>", conf.Nickname, conf.From)
	} else {
		e.From = conf.From
	}
	e.To = to
	e.Subject = subject
	e.HTML = []byte(body)

	hostAddr := fmt.Sprintf("%s:%d", conf.Host, conf.Port)
	if conf.IsSSL {
		return e.SendWithTLS(hostAddr, auth, &tls.Config{ServerName: conf.Host})
	}
	return e.Send(hostAddr, auth)
}
