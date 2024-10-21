package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const UserAgent = "codebox"

func NPMLogin(serverEndpoint string, username string, password string) (token string, err error) {
	requestBody := map[string]interface{}{
		"identity": username,
		"secret":   password,
	}
	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request body")
	}
	requestBodyReader := bytes.NewReader(requestBodyBytes)
	tokenUrl := fmt.Sprintf("%s/api/tokens", serverEndpoint)
	client := &http.Client{}
	req, _ := http.NewRequest("POST", tokenUrl, requestBodyReader)
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", "application/json")
	res, err := client.Do(req)

	if err != nil {
		return "", fmt.Errorf("failed to retrieve npm token, %s", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to retrieve npm token, recived %d from npm", res.StatusCode)
	}

	// unmarshal response
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve npm token, cannot read response, %s", err)
	}

	var bodyMap map[string]interface{}
	err = json.Unmarshal(body, &bodyMap)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve npm token, response is not a json string, %s", err)
	}

	token, ok := bodyMap["token"].(string)
	if !ok {
		return "", fmt.Errorf("failed to retrieve npm token, missing 'token' key in server response")
	}
	return token, nil
}
