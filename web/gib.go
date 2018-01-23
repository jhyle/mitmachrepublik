package mmr

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const (
	gibBasePath     = "https://www.gratis-in-berlin.de/"
	gibLoginPage    = gibBasePath + "login"
	gibLoginHandler = gibLoginPage + "?task=user.login"
	gibCreatePage   = gibBasePath + "tipp-eintragen"
	gibEditPage     = gibBasePath + "component/flexicontent/?view=item&task=edit&id=%s"
	gibEventsPage   = gibBasePath + "meine-tipps"

	GIB_CAT_SIGHTS         = 10
	GIB_CAT_TOURS          = 11
	GIB_CAT_VIEWS          = 12
	GIB_CAT_FESTIVALS      = 13
	GIB_CAT_FIREWORKS      = 14
	GIB_CAT_MUSEONS        = 21
	GIB_CAT_SHOWS          = 18
	GIB_CAT_MEETUPS        = 15
	GIB_CAT_INTERNATIONAL  = 17
	GIB_CAT_CABARET        = 19
	GIB_CAT_CINEMA         = 20
	GIB_CAT_MUSIC          = 22
	GIB_CAT_TROVE          = 26
	GIB_CAT_LISTENING      = 23
	GIB_CAT_YOUTH          = 25
	GIB_CAT_KNOWLEDGE      = 24
	GIB_CAT_WELLNESS       = 27
	GIB_CAT_OUTDOOR        = 28
	GIB_CAT_GAMES          = 29
	GIB_CAT_SPORTS         = 30
	GIB_CAT_PARTIES        = 31
	GIB_CAT_DANCING        = 32
	GIB_CAT_PUBLIC_VIEWING = 33
	GIB_CAT_BEACH_BARS     = 16
	GIB_CAT_ART            = 36
)

var (
	hiddenExp       = regexp.MustCompile(`type="hidden"[^>]+name="([^\"]+)"[^>]+value="([^\"]+)"`)
	formExp         = regexp.MustCompile(`<form[^>]+action="([^\"]+)"`)
	selectedCatExp  = regexp.MustCompile(`(?s)<select[^>]+name="jform\[catid\]".+?<option[^>]+value="([^\"]+)"[^>]+selected`)
	idInUrlExp      = regexp.MustCompile(`/(\d+)-[^/]+$`)
	paragraphEndExp = regexp.MustCompile(`</p>`)

	category2GibCat map[int]int = map[int]int{
		7:  GIB_CAT_MEETUPS,
		8:  GIB_CAT_SPORTS,
		9:  GIB_CAT_OUTDOOR,
		10: GIB_CAT_ART,
		11: GIB_CAT_KNOWLEDGE,
		12: GIB_CAT_MEETUPS,
		13: GIB_CAT_MEETUPS,
		14: GIB_CAT_MEETUPS,
		15: GIB_CAT_MEETUPS,
		16: GIB_CAT_TROVE,
		17: GIB_CAT_TROVE,
		18: GIB_CAT_OUTDOOR,
		22: GIB_CAT_YOUTH,
		23: GIB_CAT_MEETUPS,
		24: GIB_CAT_WELLNESS,
		25: GIB_CAT_TROVE,
	}
)

type GibClient struct {
	httpClient http.Client
	hostname   string
}

func (gib *GibClient) NewGetRequest(url string) (*http.Request, error) {

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating GET request for %s", url)
	}

	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), phantomUserAgent)
	return req, nil
}

func (gib *GibClient) NewPostRequest(url string, body io.Reader) (*http.Request, error) {

	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating POST request for %s", url)
	}

	req.Header.Set(http.CanonicalHeaderKey("User-Agent"), phantomUserAgent)
	return req, nil
}

func NewGibClient(hostname, user, password string) (*GibClient, error) {

	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, errors.Wrap(err, "error creating cookie jar")
	}
	client := &GibClient{hostname: hostname, httpClient: http.Client{Jar: cookieJar}}
	client.httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	req, err := client.NewGetRequest(gibLoginPage)
	if err != nil {
		return nil, errors.Wrap(err, "error creating Gratis in Berlin login request")
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "error loading Gratis in Berlin login page")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.Errorf("error loading Gratis in Berlin login page, %d", resp.StatusCode)
	}

	cookiesUrl, err := url.Parse(gibBasePath)
	if err != nil {
		return nil, errors.Wrap(err, "error parsing Gratis in Berlin base path")
	}
	client.httpClient.Jar.SetCookies(cookiesUrl, resp.Cookies())

	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading Gratis in Berlin login page")
	}

	fields := url.Values{
		"username": []string{user},
		"password": []string{password},
	}
	for _, field := range hiddenExp.FindAllSubmatch(doc, -1) {
		name := string(field[1])
		if name != "option" && name != "task" {
			fields[name] = []string{string(field[2])}
		}
	}

	post, err := client.NewPostRequest(gibLoginHandler, strings.NewReader(fields.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "error creating Gratis in Berlin login request")
	}
	post.Header.Set(http.CanonicalHeaderKey("Content-Type"), "application/x-www-form-urlencoded")

	resp, err = client.httpClient.Do(post)
	if err != nil {
		return nil, errors.Wrap(err, "error posting Gratis in Berlin login")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		return nil, errors.Errorf("error posting Gratis in Berlin login page, %d", resp.StatusCode)
	}

	redirect, err := resp.Location()
	if err != nil {
		return nil, errors.Wrap(err, "error reading Gratis in Berlin login response")
	}

	if strings.HasSuffix(redirect.String(), "/login") {
		return nil, errors.New("failed to login to Gratis in Berlin")
	}

	client.httpClient.Jar.SetCookies(cookiesUrl, resp.Cookies())
	return client, nil
}

func (client *GibClient) PostEvent(event *Event) (string, error) {

	if event.Addr.City != "Berlin" {
		return "", nil
	}

	var err error
	var req *http.Request
	if event.GiBId != "" {
		req, err = client.NewGetRequest(fmt.Sprintf(gibEditPage, event.GiBId))
		if err != nil {
			return "", errors.Wrap(err, "error creating Gratis in Berlin edit page request")
		}
	} else {
		req, err = client.NewGetRequest(gibCreatePage)
		if err != nil {
			return "", errors.Wrap(err, "error creating Gratis in Berlin create page request")
		}
	}

	resp, err := client.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error loading Gratis in Berlin event form")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.Errorf("error loading Gratis in Berlin event form, %d", resp.StatusCode)
	}

	doc, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrap(err, "error reading Gratis in Berlin event form")
	}

	form := formExp.FindSubmatch(doc)
	action := string(form[1])

	var category string
	selectedCat := selectedCatExp.FindSubmatch(doc)
	if len(selectedCat) > 0 {
		category = string(selectedCat[1])
	} else {
		if len(event.Categories) > 0 && category2GibCat[event.Categories[0]] > 0 {
			category = fmt.Sprintf("%d", category2GibCat[event.Categories[0]])
		} else {
			category = fmt.Sprintf("%d", GIB_CAT_TROVE)
		}
	}

	fields := map[string]string{}
	for _, field := range hiddenExp.FindAllSubmatch(doc, -1) {
		fields[string(field[1])] = string(field[2])
	}
	fields["task"] = "save_a_preview"
	fields["referer"] = gibEventsPage

	fields["jform[title]"] = event.Title
	fields["jform[catid]"] = category
	descr := string(event.HtmlDescription())
	descr = paragraphEndExp.ReplaceAllString(descr, "</p><br />")
	fields["jform[text]"] = descr

	fields["custom[tip_url][]"] = client.hostname + event.Url()
	fields["__fcfld_valcnt__[tip_image]"] = "0"
	fields["jform[rules][]"] = ""
	fields["custom[limits]"] = ""

	fields["custom[tip_street][]"] = event.Addr.Street
	fields["custom[tip_streetnr][]"] = ""
	fields["custom[tip_postalcode][]"] = event.Addr.Pcode
	fields["custom[tip_city]"] = event.Addr.City

	fields["custom[date_from][]"] = iso8601DateFormat(event.Start)
	fields["custom[date_till][]"] = iso8601DateFormat(event.End)
	fields["custom[time_hour]"] = fmt.Sprintf("%02d", event.Start.Hour())
	fields["custom[time_minutes]"] = fmt.Sprintf("%02d", event.Start.Minute())

	fields["custom[field31]"] = "0"
	fields["custom[turnus_month]"] = ""
	fields["custom[turnus_day]"] = ""

	if event.Recurrency != NoRecurrence {
		fields["custom[field31]"] = "1"
	}
	if event.Recurrency == Weekly {
		for _, weekday := range event.Weekly.Weekdays {
			fields["custom[turnus_day]"] = weekdayLong[weekday]
			break
		}
		if event.Weekly.Interval == 1 {
			fields["custom[turnus_month]"] = "Jeden"
		} else if event.Weekly.Interval == 2 {
			fields["custom[turnus_month]"] = "Alle 14 Tage"
		}
	} else if event.Recurrency == Monthly {
		fields["custom[turnus_day]"] = weekdayLong[event.Monthly.Weekday]
		if event.Monthly.Week == FirstWeek {
			fields["custom[turnus_month]"] = "Jeden 1."
		} else if event.Monthly.Week == SecondWeek {
			fields["custom[turnus_month]"] = "Jeden 2."
		} else if event.Monthly.Week == ThirdWeek {
			fields["custom[turnus_month]"] = "Jeden 3."
		} else if event.Monthly.Week == FourthWeek {
			fields["custom[turnus_month]"] = "Jeden 4."
		} else if event.Monthly.Week == LastWeek {
			fields["custom[turnus_month]"] = "Jeden letzten"
		}
	}

	var data bytes.Buffer
	w := multipart.NewWriter(&data)
	for key, value := range fields {
		err := w.WriteField(key, value)
		if err != nil {
			return "", errors.Wrap(err, "error creating Gratis in Berlin multipart form")
		}
	}
	w.Close()

	req, err = client.NewPostRequest(action, &data)
	if err != nil {
		return "", errors.Wrap(err, "error creating Gratis in Berlin post event request")
	}
	req.Header.Set(http.CanonicalHeaderKey("Content-Type"), w.FormDataContentType())

	resp, err = client.httpClient.Do(req)
	if err != nil {
		return "", errors.Wrap(err, "error sending Gratis in Berlin post event request")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusSeeOther {
		return "", errors.Errorf("error sending Gratis in Berlin post event request, %d", resp.StatusCode)
	}

	redirect, err := resp.Location()
	if err != nil {
		return "", errors.Wrap(err, "error reading Gratis in Berlin post event response")
	}

	if !strings.HasSuffix(redirect.String(), "preview=1") {
		return "", errors.New("failed to post to Gratis in Berlin")
	}

	id := idInUrlExp.FindStringSubmatch(redirect.String())
	return id[1], nil
}

func (client *GibClient) DeletePost(event *Event) error {

	deletedEvent := *event

	deletedEvent.Addr.City = "Berlin"
	deletedEvent.Start = time.Date(2006, time.January, 2, 15, 4, 5, 0, time.Local)
	deletedEvent.Recurrency = NoRecurrence

	_, err := client.PostEvent(&deletedEvent)
	if err != nil {
		return errors.Wrap(err, "error deleting Gratis in Berlin post")
	}

	return nil
}
