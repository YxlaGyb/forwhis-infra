package utils

import (
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"sync"
	"time"

	"forestwhisperGo/apps/pkg/utils/configutil"

	"github.com/pelletier/go-toml/v2"
)

var (
	mailCfg     map[string]map[string]interface{}
	mailCfgOnce sync.Once
)

func loadMailConfig() {
	mailCfgOnce.Do(func() {
		data, err := os.ReadFile("config.toml")
		if err != nil {
			mailCfg = make(map[string]map[string]interface{})
			return
		}
		if err := toml.Unmarshal(data, &mailCfg); err != nil {
			mailCfg = make(map[string]map[string]interface{})
			return
		}
	})
}

func mailStr(key string) string {
	loadMailConfig()
	if s, ok := mailCfg["smtp"]; ok {
		return configutil.GetString(s, key)
	}
	return ""
}

func mailInt(key string) int {
	loadMailConfig()
	if s, ok := mailCfg["smtp"]; ok {
		return configutil.GetInt(s, key)
	}
	return 0
}

func mailBool(key string) bool {
	loadMailConfig()
	if s, ok := mailCfg["smtp"]; ok {
		return configutil.GetBool(s, key)
	}
	return false
}

func SendMail(subject string, receivers []string, html string) bool {
	host, port := mailStr("MAIL_SERVER"), mailInt("MAIL_PORT")
	user, pass := mailStr("MAIL_USERNAME"), mailStr("MAIL_PASSWORD")
	name := mailStr("MAIL_SENDER_NAME")
	if name == "" {
		name = "ForestWhisper"
	}

	if host == "" || port == 0 || user == "" || pass == "" {
		return false
	}

	from := fmt.Sprintf("%s <%s>", name, user)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, strings.Join(receivers, ","), subject, html)

	auth := smtp.PlainAuth("", user, pass, host)

	addr := fmt.Sprintf("%s:%d", host, port)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}

	dialer := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	fmt.Printf("[邮件] 正在连接 %s ...\n", addr)
	conn, err := tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
	if err != nil {
		return false
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		return false
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		return false
	}

	if err = client.Mail(user); err != nil {
		return false
	}

	for _, r := range receivers {
		if err = client.Rcpt(r); err != nil {
			return false
		}
	}

	w, err := client.Data()
	if err != nil {
		return false
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		return false
	}

	if err = w.Close(); err != nil {
		return false
	}

	if err = client.Quit(); err != nil {
		return false
	}

	return true
}
