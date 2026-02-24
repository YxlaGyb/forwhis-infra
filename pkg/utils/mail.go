package utils

import (
	"crypto/tls"
	"fmt"
	"net/smtp"
	"os"
	"strings"
	"sync"

	"forestwhisperGo/pkg/utils/configutil"

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
	if s, ok := mailCfg["smtp"]; ok {
		return configutil.GetString(s, key)
	}
	return ""
}

func mailInt(key string) int {
	if s, ok := mailCfg["smtp"]; ok {
		return configutil.GetInt(s, key)
	}
	return 0
}

func mailBool(key string) bool {
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

	from := fmt.Sprintf("%s <%s>", name, user)
	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, strings.Join(receivers, ","), subject, html)

	auth := smtp.PlainAuth("", user, pass, host)

	addr := fmt.Sprintf("%s:%d", host, port)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: mailBool("MAIL_SKIP_TLS_VERIFY"),
		ServerName:         host,
	}

	conn, err := tls.Dial("tcp", addr, tlsConfig)
	if err != nil {
		fmt.Printf("邮件连接失败: %v\n", err)
		return false
	}
	defer conn.Close()

	client, err := smtp.NewClient(conn, host)
	if err != nil {
		fmt.Printf("创建SMTP客户端失败: %v\n", err)
		return false
	}
	defer client.Close()

	if err = client.Auth(auth); err != nil {
		fmt.Printf("SMTP认证失败: %v\n", err)
		return false
	}

	if err = client.Mail(user); err != nil {
		fmt.Printf("设置发件人失败: %v\n", err)
		return false
	}

	for _, r := range receivers {
		if err = client.Rcpt(r); err != nil {
			fmt.Printf("设置收件人失败: %v\n", err)
			return false
		}
	}

	w, err := client.Data()
	if err != nil {
		fmt.Printf("准备邮件数据失败: %v\n", err)
		return false
	}

	_, err = w.Write([]byte(msg))
	if err != nil {
		fmt.Printf("写入邮件内容失败: %v\n", err)
		return false
	}

	if err = w.Close(); err != nil {
		fmt.Printf("关闭邮件写入失败: %v\n", err)
		return false
	}

	if err = client.Quit(); err != nil {
		fmt.Printf("退出SMTP连接失败: %v\n", err)
		return false
	}

	return true
}
