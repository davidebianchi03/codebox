package emails

import (
	"bytes"
	"html/template"
	"path"
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

	baseTemplatePath := path.Join(
		config.Environment.TemplatesFolder,
		"emails",
		"html",
	)

	tmpl, err := template.ParseFiles(
		path.Join(baseTemplatePath, "email_base.html"),
		path.Join(baseTemplatePath, templateName),
	)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	if err := tmpl.ExecuteTemplate(buf, "email_base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}

/*
RenderTextEmailTemplate: renders an text email template with the provided data.

Parameters:
- templateName: the name of the template to render (e.g., "email_verify_address").
- data: a map containing the data to be injected into the template.
*/
func RenderTextEmailTemplate(templateName string, data map[string]any) (string, error) {
	data["year"] = time.Now().Year()
	data["serverURL"] = config.Environment.ExternalUrl

	baseTemplatePath := path.Join(
		config.Environment.TemplatesFolder,
		"emails",
		"text",
	)

	tmpl, err := template.ParseFiles(
		path.Join(baseTemplatePath, "email_base.txt"),
		path.Join(baseTemplatePath, templateName),
	)
	if err != nil {
		return "", err
	}

	buf := bytes.NewBufferString("")

	if err := tmpl.ExecuteTemplate(buf, "email_base", data); err != nil {
		return "", err
	}

	return buf.String(), nil
}
