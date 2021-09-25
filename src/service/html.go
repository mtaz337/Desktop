package service

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"html/template"
)

var (
	emailTemples = map[string]string{
		"OTP": `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>AirBringr</title>
	</head>
	<body>
		This is your OTP {{.otp}}
	</body>
</html>`,
	}
)

func GenerateTpl(tplCode string, data interface{}) (string, error) {
	tmpl := emailTemples[tplCode]
	t, err := template.New("email").Parse(tmpl)

	if err != nil {
		log.Error(err.Error())
		return "", err
	}

	var html bytes.Buffer
	err = t.Execute(&html, data)
	return html.String(), nil
}
