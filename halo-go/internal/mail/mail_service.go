package mail

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"strings"
)

type MailService struct {
	Host     string
	Port     int
	Username string
	Password string
	From     string
	FromName string
}

type MailConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	From     string `json:"from"`
	FromName string `json:"fromName"`
}

func NewMailService(cfg *MailConfig) *MailService {
	return &MailService{
		Host:     cfg.Host,
		Port:     cfg.Port,
		Username: cfg.Username,
		Password: cfg.Password,
		From:     cfg.From,
		FromName: cfg.FromName,
	}
}

type Email struct {
	To      string
	Subject string
	Body    string
	IsHTML  bool
}

func (m *MailService) Send(email *Email) error {
	if m.Host == "" || m.Port == 0 {
		return fmt.Errorf("邮件服务未配置")
	}

	auth := smtp.PlainAuth("", m.Username, m.Password, m.Host)

	from := fmt.Sprintf("%s <%s>", m.FromName, m.From)
	to := email.To

	headers := make(map[string]string)
	headers["From"] = from
	headers["To"] = to
	headers["Subject"] = email.Subject
	headers["MIME-Version"] = "1.0"
	if email.IsHTML {
		headers["Content-Type"] = `text/html; charset="UTF-8"`
	} else {
		headers["Content-Type"] = `text/plain; charset="UTF-8"`
	}

	message := ""
	for k, v := range headers {
		message += fmt.Sprintf("%s: %s\r\n", k, v)
	}
	message += "\r\n" + email.Body

	addr := fmt.Sprintf("%s:%d", m.Host, m.Port)

	var err error
	if m.Port == 465 {
		err = m.sendWithTLS(addr, auth, from, to, []byte(message))
	} else {
		err = smtp.SendMail(addr, auth, m.From, []string{to}, []byte(message))
	}

	return err
}

func (m *MailService) sendWithTLS(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	conn, err := tls.Dial("tcp", addr, &tls.Config{
		InsecureSkipVerify: false,
		ServerName:         strings.Split(addr, ":")[0],
	})
	if err != nil {
		return err
	}

	client, err := smtp.NewClient(conn, addr)
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Auth(a); err != nil {
		return err
	}

	if err = client.Mail(from); err != nil {
		return err
	}

	for _, addr := range to {
		if err = client.Rcpt(addr); err != nil {
			return err
		}
	}

	w, err := client.Data()
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

	return client.Quit()
}

func (m *MailService) SendCommentNotification(toEmail, postTitle, author, content string) error {
	body := fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif;">
<h2>新评论通知</h2>
<p>您发布的文章 <strong>%s</strong> 收到了新评论：</p>
<div style="background-color: #f5f5f5; padding: 15px; border-radius: 5px; margin: 10px 0;">
<p><strong>评论者：</strong>%s</p>
<p><strong>内容：</strong></p>
<p>%s</p>
</div>
<p>请登录后台查看：<a href="/console">管理后台</a></p>
</body>
</html>
`, postTitle, author, content)

	email := &Email{
		To:      toEmail,
		Subject: fmt.Sprintf("[Halo] 文章《%s》收到新评论", postTitle),
		Body:    body,
		IsHTML:  true,
	}

	return m.Send(email)
}

func (m *MailService) SendSystemAlert(subject, message string) error {
	email := &Email{
		To:      m.From,
		Subject: fmt.Sprintf("[Halo 系统告警] %s", subject),
		Body: fmt.Sprintf(`
<html>
<body style="font-family: Arial, sans-serif;">
<h2>系统告警</h2>
<p><strong>时间：</strong>%s</p>
<p><strong>告警信息：</strong></p>
<pre style="background-color: #f5f5f5; padding: 15px;">%s</pre>
</body>
</html>
`, subject, message),
		IsHTML: true,
	}

	return m.Send(email)
}
