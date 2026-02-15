package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"sync"

	"github.com/pelletier/go-toml/v2"
)

var (
	mailCfg     map[string]map[string]interface{}
	mailCfgOnce sync.Once
)

func mailVal(key string) interface{} {
	mailCfgOnce.Do(func() {
		data, _ := os.ReadFile("config.toml")
		toml.Unmarshal(data, &mailCfg)
	})
	if s, ok := mailCfg["smtp"]; ok {
		return s[key]
	}
	return nil
}

func mailStr(key string) string {
	if v, ok := mailVal(key).(string); ok {
		return v
	}
	return ""
}

func mailInt(key string) int {
	switch n := mailVal(key).(type) {
	case int:
		return n
	case int64:
		return int(n)
	case float64:
		return int(n)
	}
	return 0
}

func SendMail(subject string, receivers []string, html string) bool {
	host, port := mailStr("MAIL_SERVER"), mailInt("MAIL_PORT")
	user, pass := mailStr("MAIL_USERNAME"), mailStr("MAIL_PASSWORD")
	name := mailStr("MAIL_SENDER_NAME")
	if name == "" {
		name = "ForestWhisper"
	}

	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("From: %s <%s>\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		name, user, receivers[0], subject, html))

	auth := smtp.PlainAuth("", user, pass, host)
	conn, err := tls.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &tls.Config{InsecureSkipVerify: true, ServerName: host})
	if err != nil {
		return false
	}
	defer conn.Close()

	client, _ := smtp.NewClient(conn, host)
	defer client.Close()
	client.Auth(auth)
	client.Mail(user)
	for _, r := range receivers {
		client.Rcpt(r)
	}
	w, _ := client.Data()
	w.Write(buf.Bytes())
	w.Close()
	client.Quit()
	return true
}
