package gmailnator

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func (session *Session) Init(proxy *string) error {
	session.Proxy = proxy

	req, err := http.NewRequest("GET", "https://www.emailnator.com", nil)
	if err != nil {
		return err
	}

	client := &http.Client{}

	if session.Proxy != nil {
		proxyUrl, err := url.Parse("http://" + *session.Proxy)
		if err != nil {
			return err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	session.XsrfToken = res.Cookies()[0].Value
	session.GmailnatorSession = res.Cookies()[1].Value

	return nil
}

func (session Session) IsAlive() (bool, error) {
	req, err := http.NewRequest("POST", "https://www.emailnator.com/message-list", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Cookie", "XSRF-TOKEN=" + session.XsrfToken + "; gmailnator_session=" + session.GmailnatorSession)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Xsrf-Token", strings.ReplaceAll(session.XsrfToken, "%3D", "="))

	client := &http.Client{}

	if session.Proxy != nil {
		proxyUrl, err := url.Parse("http://" + *session.Proxy)
		if err != nil {
			return false, err
		}
		client = &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}
	}

	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if strings.Contains(string(body), "Page Expired") {
		return false, nil
	}

	return true, nil
}
