package service

import (
	"bytes"
	"html/template"

	log "github.com/sirupsen/logrus"
)

var (
	emailTemples = map[string]string{
		"otp": `
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

		"signup_verification": `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>AirBringr Email Verification</title>
	</head>
	<body>
		Hello,
		<br>
		To verify your account please click the link below. <strong>This verification link will be valid for 3 days.</strong>
		<br>
		{{.link}}
	</body>
</html>
`,
		"password_reset": `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>AirBringr Password Reset</title>
	</head>
	<body>
		Hello,
		<br>
		To reset your password please click the link below.
		<br>
		{{.link}}
	</body>
</html>
`,

		"password_reset_confirmation": `
<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>AirBringr Password Reset Confirmation</title>
	</head>
	<body>
		Hello,
		<br>
		Your password has been updated successfully.
	</body>
</html>
`,
	}
)

func GenerateTpl(tplCode string, data interface{}) (htmlStr string, err error) {
	var t *template.Template
	tmpl := emailTemples[tplCode]
	if t, err = template.New("email").Parse(tmpl); err != nil {
		log.Error(err.Error())
		return "", err
	}

	var html bytes.Buffer
	if err = t.Execute(&html, data); err != nil {
		log.Error(err.Error())
		return "", err
	}
	return html.String(), nil
}
