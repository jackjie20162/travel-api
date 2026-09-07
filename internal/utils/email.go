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
		// QQ 邮箱 465 端口使用隐式 TLS（SMTPS）：先建立 TLS 连接，再走 SMTP 协议
		tlsConn, err := tls.Dial("tcp", hostAddr, &tls.Config{ServerName: conf.Host})
		if err != nil {
			return fmt.Errorf("TLS 连接失败: %w", err)
		}
		defer tlsConn.Close()

		client, err := smtp.NewClient(tlsConn, conf.Host)
		if err != nil {
			return fmt.Errorf("SMTP 客户端创建失败: %w", err)
		}
		defer client.Close()

		auth := smtp.PlainAuth("", conf.From, conf.Secret, conf.Host)
		if err := client.Auth(auth); err != nil {
			return fmt.Errorf("SMTP 认证失败: %w", err)
		}
		if err := client.Mail(conf.From); err != nil {
			return fmt.Errorf("MAIL FROM 失败: %w", err)
		}
		for _, addr := range to {
			if err := client.Rcpt(addr); err != nil {
				return fmt.Errorf("RCPT TO 失败: %w", err)
			}
		}
		w, err := client.Data()
		if err != nil {
			return fmt.Errorf("DATA 失败: %w", err)
		}
		msg, err := e.Bytes()
		if err != nil {
			return fmt.Errorf("生成邮件内容失败: %w", err)
		}
		_, err = w.Write(msg)
		if err != nil {
			return fmt.Errorf("写入邮件内容失败: %w", err)
		}
		err = w.Close()
		if err != nil {
			return fmt.Errorf("关闭邮件数据失败: %w", err)
		}
		return client.Quit()
	}

	auth := smtp.PlainAuth("", conf.From, conf.Secret, conf.Host)
	return e.Send(hostAddr, auth)
}
