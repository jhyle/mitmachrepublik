package mmr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/benbjohnson/phantomjs"
	"github.com/pkg/errors"
)

const (
	phantomUserAgent = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:57.0) Gecko/20100101 Firefox/57.0"
	fbLoginPage      = "https://www.facebook.com/v2.11/dialog/oauth?client_id=%s&redirect_uri=https%%3A%%2F%%2Fwww.mitmachrepublik.de%%2Foauth"
	fbGetAccessToken = "https://graph.facebook.com/v2.11/oauth/access_token?client_id=%s&redirect_uri=https%%3A%%2F%%2Fwww.mitmachrepublik.de%%2Foauth&client_secret=%s&code=%s"
	fbGetMyAccounts  = "https://graph.facebook.com/v2.11/me/accounts?access_token=%s"
	fbPageName       = "Mitmach-Republik e.V."
	fbPageFeed       = "https://graph.facebook.com/v2.11/mitmachrepublik/feed?access_token=%s"
	fbPagePost       = "https://graph.facebook.com/v2.11/%s?access_token=%s"
)

type FacebookClient struct {
	accessToken string
	hostname    string
}

func NewFacebookClient(hostname, appId, appSecret, user, password string) (*FacebookClient, error) {

	p := phantomjs.NewProcess()
	err := p.Open()
	if err != nil {
		return nil, errors.Wrap(err, "err creating browser")
	}
	defer p.Close()

	page, err := p.CreateWebPage()
	if err != nil {
		return nil, errors.Wrap(err, "err creating web page")
	}

	settings, err := page.Settings()
	if err != nil {
		return nil, errors.Wrap(err, "err loading settings")
	}
	settings.UserAgent = phantomUserAgent

	err = page.Open(fmt.Sprintf(fbLoginPage, appId))
	if err != nil {
		return nil, errors.Wrap(err, "err loading login page")
	}

	time.Sleep(3 * time.Second)

	_, err = page.Content()
	if err != nil {
		return nil, errors.Wrap(err, "err loading login page")
	}

	_, err = page.Evaluate(`
		function() {
			document.body.querySelector('input[name="email"]').setAttribute('value', '` + user + `');
			document.body.querySelector('input[name="pass"]').setAttribute('value', '` + password + `');

		    var ev = document.createEvent("MouseEvent");
		    ev.initMouseEvent(
		        "click",
		        true /* bubble */, true /* cancelable */,
		        window, null,
		        0, 0, 0, 0, /* coordinates */
		        false, false, false, false, /* modifier keys */
		        0 /*left*/, null
		    );
		    document.body.querySelector('input[name="login"]').dispatchEvent(ev);
		}
	`)
	if err != nil {
		return nil, errors.Wrap(err, "err using login form")
	}

	time.Sleep(3 * time.Second)

	oAuthUrl, err := page.URL()
	if err != nil {
		return nil, errors.Wrap(err, "err getting OAuth url")
	}

	parsedUrl, err := url.Parse(oAuthUrl)
	if err != nil {
		return nil, errors.Wrap(err, "err parsing OAuth url")
	}

	code := parsedUrl.Query().Get("code")
	resp, err := http.Get(fmt.Sprintf(fbGetAccessToken, appId, appSecret, code))
	if err != nil {
		return nil, errors.Wrap(err, "error getting access token")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error reading access token, status code %d", resp.StatusCode)
	}

	var result map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, errors.Wrap(err, "error reading get token response")
	}

	return &FacebookClient{accessToken: result["access_token"].(string), hostname: hostname}, nil
}

type fbMyAccounts struct {
	Data []struct {
		AccessToken string   `json:"access_token"`
		Category    string   `json:"category"`
		Name        string   `json:"name"`
		ID          string   `json:"id"`
		Perms       []string `json:"perms"`
	} `json:"data"`
	Paging struct {
		Cursors struct {
			Before string `json:"before"`
			After  string `json:"after"`
		} `json:"cursors"`
	} `json:"paging"`
}

func (fb *FacebookClient) pageToken() (string, error) {

	resp, err := http.Get(fmt.Sprintf(fbGetMyAccounts, fb.accessToken))
	if err != nil {
		return "", errors.Wrap(err, "error reading user accounts")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("error reading user accounts, status code %d", resp.StatusCode)
	}

	var accounts fbMyAccounts
	err = json.NewDecoder(resp.Body).Decode(&accounts)
	if err != nil {
		return "", errors.Wrap(err, "error decoding accounts response")
	}

	var pageToken string
	for _, account := range accounts.Data {
		if account.Name == fbPageName {
			pageToken = account.AccessToken
			break
		}
	}

	if pageToken == "" {
		return "", errors.New("page token not found")
	}

	return pageToken, nil
}

type fbPost struct {
	Link string `json:"link"`
}

type fbId struct {
	Id string `json:"id"`
}

func (fb *FacebookClient) PostEvent(event *Event) (string, error) {

	pageToken, err := fb.pageToken()
	if err != nil {
		return "", errors.Wrap(err, "error retrieving page token")
	}

	data, err := json.Marshal(&fbPost{Link: fb.hostname + event.Url()})
	if err != nil {
		return "", errors.Wrap(err, "error encoding post")
	}

	resp, err := http.Post(fmt.Sprintf(fbPageFeed, pageToken), "application/json", bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "error posting to feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("error posting to feed, status code %d", resp.StatusCode)
	}

	var result fbId
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", errors.Wrap(err, "error reading post response")
	}

	return result.Id, nil
}

func (fb *FacebookClient) DeletePost(id string) error {

	pageToken, err := fb.pageToken()
	if err != nil {
		return errors.Wrap(err, "error retrieving page token")
	}

	req, err := http.NewRequest(http.MethodDelete, fmt.Sprintf(fbPagePost, id, pageToken), nil)
	if err != nil {
		return errors.Wrap(err, "error creating delete request")
	}

	resp, err := (&http.Client{}).Do(req)
	if err != nil {
		return errors.Wrap(err, "error posting to feed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("error deleting posting, status code %d", resp.StatusCode)
	}

	return nil
}
