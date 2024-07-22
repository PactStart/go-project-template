package email

import (
	"bytes"
	"context"
	"github.com/pkg/errors"
	"html/template"
	"net/smtp"
	"orderin-server/pkg/common/log"
	"orderin-server/pkg/common/utils"
	"path/filepath"
	"strings"
)

type QQEmail struct {
	Config *EmailConfig
}

func (e *QQEmail) Setup(config *EmailConfig) error {
	e.Config = config
	return nil
}

func (e *QQEmail) Send(request EmailRequest) error {
	var body string
	var contentType string
	var err error
	if request.Content != "" {
		body = request.Content
		contentType = "text/plain"
	} else if request.TemplateFileName != "" {
		// 读取模板文件
		templatePath := filepath.Join("templates", request.TemplateFileName)
		tmpl, err := template.ParseFiles(templatePath)
		if err != nil {
			log.ZError(context.Background(), "Error read template:", err)
			return err
		}
		// 渲染模板
		buf := &bytes.Buffer{}
		err = tmpl.Execute(buf, request.Params)
		if err != nil {
			log.ZError(context.Background(), "Error parsing template:", err)
			return err
		}
		body = buf.String()
		contentType = "text/html"
	}
	if body != "" {
		//发送
		message := []byte("To: " + strings.Join(request.To, ", ") + "\r\n" +
			"Subject: " + request.Subject + "\r\n" +
			"Content-Type: " + contentType + "; charset=UTF-8\r\n" +
			"From: " + e.Config.From + "\r\n" + "\r\n" +
			body)

		auth := smtp.PlainAuth("", e.Config.Account, e.Config.Password, e.Config.SmtpHost)
		err = smtp.SendMail(e.Config.SmtpHost+":"+utils.IntToString(e.Config.SmtpPort), auth, e.Config.Account, request.To, message)
		if err != nil {
			log.ZError(context.Background(), "Error send email:", err)
			return err
		}
		return nil
	} else {
		return errors.New("Cannot find available email content")
	}
}
