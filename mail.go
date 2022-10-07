package gmailnator

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"errors"
	"bytes"
)

func (session Session) RetrieveMail(email string) ([]Email, error) {
	var emails []Email

	values := map[string]string{"email": email}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return emails, err
	}

	req, err := http.NewRequest("POST", "https://www.emailnator.com/message-list", bytes.NewBuffer(jsonValue))
	if err != nil {
		return emails, err
	}

	req.Header.Set("Cookie", "XSRF-TOKEN=" + session.XsrfToken + "; gmailnator_session=" + session.GmailnatorSession)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", strings.ReplaceAll(session.XsrfToken, "%3D", "="))

	client := &http.Client{}

	if session.Proxy != nil {
		proxyUrl, err := url.Parse("http://" + *session.Proxy)
		if err != nil {
			return emails, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	resp, err := client.Do(req)
	if err != nil {
		return emails, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return emails, &RequestErr{
			StatusCode: resp.StatusCode,
			Err:        errors.New("invalid status code"),
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return emails, err
    }

	var jsonResponse RetrieveEmailsJsonResponse

	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		return emails, err
	}
	
	emails = append(emails, jsonResponse.MessageData...)[1:]

	return emails, nil
}

func (session Session) GenerateEmailAddress() (string, error) {
	values := map[string][]string{"email": {"plusGmail", "dotGmail"}}

	jsonValue, err := json.Marshal(values)
	if err != nil {
		return "", err
	}
	
	req, err := http.NewRequest("POST", "https://www.emailnator.com/generate-email", bytes.NewBuffer(jsonValue))
	if err != nil {
		return "", err
	}

	req.Header.Set("Cookie", "XSRF-TOKEN=" + session.XsrfToken + "; gmailnator_session=" + session.GmailnatorSession)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", strings.ReplaceAll(session.XsrfToken, "%3D", "="))

	client := &http.Client{}

	if session.Proxy != nil {
		proxyUrl, err := url.Parse("http://" + *session.Proxy)
		if err != nil {
			return "", err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", &RequestErr{
			StatusCode: resp.StatusCode,
			Err:        errors.New("invalid status code"),
		}
	}

	var jsonResponse GenerateEmailJsonResponse

	body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

	err = json.Unmarshal(body, &jsonResponse)
    if err != nil {
        return "", err
    }

	return jsonResponse.Email[0], nil
}
