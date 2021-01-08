package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"
)

func NewHttpApiHandler(url, accessToken string) *HttpApiHandler {
	handler := &HttpApiHandler{
		url:         url,
		accessToken: accessToken,
	}
	return handler
}

func (handler *HttpApiHandler) Request(method string, path string, body interface{}, result interface{}) error {

	url := handler.url + path

	var client = &http.Client{Timeout: 10 * time.Second}
	var err error
	var resp *http.Response

	hasBody := body != nil && (method == "POST" || method == "PUT")

	var reader io.Reader

	var b []byte
	if hasBody {
		b, err = json.Marshal(body)
		if err != nil {
			return err
		}
		reader = bytes.NewBuffer(b)
	}

	if os.Getenv("DEBUG_HTTP_ALL") != "" || (hasBody && os.Getenv("DEBUG_HTTP") != "") {
		fmt.Println("--> request Method:", method)
		fmt.Println("--> request Url:", url)
		fmt.Printf("--> request Body: %s\n", string(b))
	}

	req, err := http.NewRequest(method, url, reader)
	if err != nil {
		return err
	}

	if hasBody {
		req.Header.Set("Content-Type", "application/json")
	}

	if handler.accessToken != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", handler.accessToken))
	}

	resp, err = client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode >= 400 {
		errorBody := &ErrorResponse{}
		err = json.Unmarshal(respBody, &errorBody)
		if err == nil {
			return fmt.Errorf("%d %s", resp.StatusCode, errorBody.Error)
		} else {
			return fmt.Errorf("Bad status code %d", resp.StatusCode)
		}
	}

	if os.Getenv("DEBUG_HTTP_ALL") != "" || (hasBody && os.Getenv("DEBUG_HTTP") != "") {
		fmt.Println("--> response Status:", resp.Status)
		fmt.Println("--> response Headers:", resp.Header)
		fmt.Println("--> response Body:", string(respBody))
	}

	if result != nil {
		bodyString := string(respBody)
		err = json.Unmarshal(respBody, &result)
		if err != nil {
			log.Printf("Error decoding response '%s'", bodyString)
		}
	}
	return err
}
