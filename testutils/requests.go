package testutils

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"gitlab.com/codebox4073715/codebox/db/models"
)

func AuthenticateHttpRequest(t *testing.T, req *http.Request, user models.User) {
	// create a fake token for the user
	token, err := models.CreateToken(user, time.Duration(time.Hour*24*20))
	if err != nil {
		t.Errorf("Failed to create token: '%s'\n", err)
		t.FailNow()
	}

	req.Header.Set("Authorization", "Bearer "+token.Token)
}

func CreateRequestWithJSONBody(t *testing.T, url, method string, body interface{}) *http.Request {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		t.Errorf("Failed to create HTTP request: '%s'\n", err)
		t.FailNow()
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		t.Errorf("Failed to marshal request body: '%s'\n", err)
		t.FailNow()
	}

	req.Body = io.NopCloser(bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	return req
}
