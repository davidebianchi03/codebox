package emails

import (
	"bytes"
	"html/template"
	"io"
	"time"

	"gitlab.com/codebox4073715/codebox/config"
)

/*
RenderHtmlEmailTemplate: renders an HTML email template with the provided data.

Parameters:
- templateName: the name of the template to render (e.g., "email_verify_address_html").
- data: a map containing the data to be injected into the template.
*/
func RenderHtmlEmailTemplate(templateName string, data map[string]any) (string, error) {
	data["year"] = time.Now().Year()
	data["serverURL"] = config.Environment.ExternalUrl

	tmpl, err := template.ParseGlob("html/emails/html/*.html")
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	err = tmpl.ExecuteTemplate(io.Discard, templateName, data)
	if err != nil {
		return "", err
	}

	err = tmpl.ExecuteTemplate(buf, "email_base_html", data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
