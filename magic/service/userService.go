package service

import "errors"

var ErrInvalidCredentials = errors.New("邮箱或密码错误喵")

type LoginResult struct {
	JID   string
	UID   int
	Name  string
	Email string
}

func LoginUser(email, password, clientIP string) (*LoginResult, error) {
	// TODO: 实现数据库验证
	if email == "test@example.com" && password == "password" {
		return &LoginResult{
			JID:   "test_jid",
			UID:   1,
			Name:  "TestUser",
			Email: email,
		}, nil
	}
	return nil, ErrInvalidCredentials
}
