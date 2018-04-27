package mmr

import (
	"bytes"
	"html"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kennygrant/sanitize"
	"github.com/pkg/errors"
	"gopkg.in/mgo.v2/bson"
)

type (
	ScrapersService struct {
		BasicService
		events      *EventService
		organizerId bson.ObjectId
	}
)

const (
	UK_SOURCE      = "umweltkalender"
	UK_LIST_REGEX  = "<li id=\"(\\d+)\" class=\"list_item\">"
	UK_LIST_URL    = "http://www.umweltkalender-berlin.de/angebote/aktuell?p=1&arten[0]=MIMA&um_km=500&zeitraum=zz"
	UK_DETAILS_URL = "http://www.umweltkalender-berlin.de/angebote/details/"

	TAG_TITLE_REGEX     = "<title[^>]*>([^<]*)</title"
	TAG_H1_REGEX        = "<h1[^>]*>(.*)</h1"
	TAG_H2_REGEX        = "<h2[^>]*>(.*)</h2"
	TAG_P_CLEAN_REGEX   = "<p[^>]*>(.*)</p>"
	DATE_FROM_REGEX     = "(\\d+)\\.(\\d+)\\.(\\d+), (\\d+):(\\d+)"
	DATE_FROM_TO_REGEX  = "(\\d+)\\.(\\d+)\\.(\\d+), (\\d+):(\\d+) - (\\d+):(\\d+)"
	POSTCODE_CITY_REGEX = "(\\d{5})\\s+(.*)"
	POSTCODE_REGEX      = "(\\d{5})"

	UK_TARGETS_REGEX  = "<p><strong>Für:</strong>([^<]*)<"
	UK_LOCATION_REGEX = "(?s)<p><strong>Ort / Start:</strong>([^<]*)<"
	UK_DIV_MAIN_REGEX = "(?s)<div class=\"green_corners_main\">(.*)<!-- end: div green_corner_main -->"
)

var (
	ANMELDUNG_ERFORDERLICH = []byte("ANMELDUNG ERFORDERLICH!")

	uk_list       *regexp.Regexp = regexp.MustCompile(UK_LIST_REGEX)
	tag_title     *regexp.Regexp = regexp.MustCompile(TAG_TITLE_REGEX)
	tag_h1        *regexp.Regexp = regexp.MustCompile(TAG_H1_REGEX)
	tag_h2        *regexp.Regexp = regexp.MustCompile(TAG_H2_REGEX)
	tag_p_clean   *regexp.Regexp = regexp.MustCompile(TAG_P_CLEAN_REGEX)
	date_from_to  *regexp.Regexp = regexp.MustCompile(DATE_FROM_TO_REGEX)
	date_from     *regexp.Regexp = regexp.MustCompile(DATE_FROM_REGEX)
	postcode_city *regexp.Regexp = regexp.MustCompile(POSTCODE_CITY_REGEX)
	postcode      *regexp.Regexp = regexp.MustCompile(POSTCODE_REGEX)
	uk_div_main   *regexp.Regexp = regexp.MustCompile(UK_DIV_MAIN_REGEX)
	uk_targets    *regexp.Regexp = regexp.MustCompile(UK_TARGETS_REGEX)
	uk_locations  *regexp.Regexp = regexp.MustCompile(UK_LOCATION_REGEX)
	zero_time     time.Time      = time.Date(0, time.Month(1), 0, 0, 0, 0, 0, time.Local)

	importedTags       = []string{"div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "table", "tbody", "tr", "td"}
	importedAttributes = []string{"title"}
)

func NewScrapersService(hour int, email *EmailAccount, events *EventService, organizerId bson.ObjectId) Service {

	return &ScrapersService{NewBasicService("NewSraperService", hour, email), events, organizerId}
}

func (service *ScrapersService) Start() {

	service.start(service.Run)
}

func (service *ScrapersService) getWebPage(url string) ([]byte, error) {

	resp, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting url: %s", url)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading body of url: %s", url)
	}

	return body, nil
}

func (service *ScrapersService) saveScraped(event *Event) error {

	oldEvent, err := service.events.LoadBySource(event.Source, event.SourceId)
	if err == nil {
		if oldEvent.Start != event.Start || oldEvent.End != event.End || oldEvent.Descr != event.Descr || oldEvent.Title != event.Title {
			oldEvent.Start = event.Start
			oldEvent.End = event.End
			oldEvent.Title = event.Title
			oldEvent.Descr = event.Descr
			err = service.events.Store(oldEvent)
			if err != nil {
				return errors.Wrapf(err, "error updating imported event: %s %s", event.Source, event.SourceId)
			}
		}
	} else {
		event.Id = bson.NewObjectId()
		event.OrganizerId = service.organizerId
		err = service.events.Store(event)
		if err != nil {
			return errors.Wrapf(err, "error storing imported event: %s %s", event.Source, event.SourceId)
		}
	}

	return nil
}

func (service *ScrapersService) makeDateTime(day, month, year, hour, minute string) (time.Time, error) {

	iDay, err := strconv.Atoi(day)
	if err != nil {
		return zero_time, err
	}

	iMonth, err := strconv.Atoi(month)
	if err != nil {
		return zero_time, err
	}

	iYear, err := strconv.Atoi(year)
	if err != nil {
		return zero_time, err
	}

	iHour, err := strconv.Atoi(hour)
	if err != nil {
		return zero_time, err
	}

	iMinute, err := strconv.Atoi(minute)
	if err != nil {
		return zero_time, err
	}

	return time.Date(iYear, time.Month(iMonth), iDay, iHour, iMinute, 0, 0, time.Local), nil
}

func (service *ScrapersService) scrapeUmweltKalenderEvent(id string) (*Event, error) {

	var event Event
	event.Source = UK_SOURCE
	event.SourceId = id
	event.Web = UK_DETAILS_URL + id

	page, err := service.getWebPage(event.Web)
	if err != nil {
		return nil, errors.Wrapf(err, "error scraping Umweltkalender page: %s", id)
	}

	title := tag_title.FindSubmatch(page)
	if title == nil {
		return nil, errors.New("could not find title at Umweltkalender #" + id)
	}

	var date [][]byte
	if date_from_to.Match(title[1]) {
		date = date_from_to.FindSubmatch(title[1])
	} else if date_from.Match(title[1]) {
		date = date_from.FindSubmatch(title[1])
	} else {
		return nil, errors.New("could not find date at Umweltkalender #" + id)
	}

	event.Start, err = service.makeDateTime(string(date[1]), string(date[2]), "20"+string(date[3]), string(date[4]), string(date[5]))
	if err != nil {
		return nil, errors.Wrapf(err, "error parsing start date at Umweltkalender page: %s", id)
	}

	if len(date) > 6 {
		event.End, err = service.makeDateTime(string(date[1]), string(date[2]), "20"+string(date[3]), string(date[6]), string(date[7]))
		if err != nil {
			return nil, errors.Wrapf(err, "error parsing end date at Umweltkalender page: %s", id)
		}
	}

	h1 := tag_h1.FindSubmatch(page)
	if title == nil {
		return nil, errors.New("could not find h1 at Umweltkalender #" + id)
	}
	event.Title = html.UnescapeString(string(h1[1]))

	main := uk_div_main.FindSubmatch(page)
	if main == nil {
		return nil, errors.New("could not find main div at Umweltkalender #" + id)
	}
	paras := tag_p_clean.FindAll(main[1], -1)
	for i := range paras {
		event.Descr += string(paras[i])
	}

	event.Descr, err = sanitize.HTMLAllowing(event.Descr, importedTags, importedAttributes)
	if err != nil {
		return nil, errors.Wrapf(err, "error sanitizing HTML of Umweltkalende page: %s", id)
	}
	event.Descr = strings.Replace(event.Descr, "»»»", "", -1)
	event.Descr += "<p class=\"small\">Text von der Veranstaltungswebseite übernommen</p>"

	location := uk_locations.FindSubmatch(page)
	if location == nil {
		return nil, errors.New("could not find location at Umweltkalender #" + id)
	}

	postcode_part := -1
	address := bytes.Split(location[1], []byte(", "))
	for i := range address {
		if postcode.Match(address[i]) {
			postcode_part = i
			break
		}
	}
	if postcode_part > 0 {
		event.Addr.Street = strings.Trim(string(address[postcode_part-1]), " ")
	}
	if postcode_part > -1 {
		district := postcode_city.FindSubmatch(address[postcode_part])
		if district != nil {
			event.Addr.Pcode = string(district[1])
			event.Addr.City = strings.Trim(string(district[2]), " ")
		} else {
			pcode := postcode.FindSubmatch(address[postcode_part])
			if pcode == nil {
				return nil, errors.New("could not find postcode or city at Umweltkalender #" + id)
			}
			event.Addr.Pcode = string(pcode[1])
			event.Addr.City = "Berlin"
		}
	}

	targets := uk_targets.FindSubmatch(page)
	if main == nil {
		return nil, errors.New("could not find targets at Umweltkalender #" + id)
	}

	event.Targets = make([]int, 0)
	for _, target := range bytes.Split(targets[1], []byte(", ")) {
		targetId, exists := TargetMap[strings.Trim(string(target), " ")]
		if exists {
			event.Targets = append(event.Targets, targetId)
		}
	}

	event.Categories = []int{CategoryMap["Leute treffen"], CategoryMap["Umweltschutz"]}
	event.Rsvp = bytes.Contains(page, ANMELDUNG_ERFORDERLICH)
	return &event, nil
}

func (service *ScrapersService) scrapeUmweltKalender() []error {

	page, err := service.getWebPage(UK_LIST_URL)
	if err != nil {
		return []error{errors.Wrapf(err, "error loading Umweltkalender event list at %s", UK_LIST_URL)}
	}

	ids := uk_list.FindAllSubmatch(page, -1)
	if ids == nil {
		return []error{errors.New("no events found at Umweltkalender")}
	}

	errs := []error{}
	for i := range ids {
		event, err := service.scrapeUmweltKalenderEvent(string(ids[i][1]))
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error reading Umweltkalender event: %s", ids[i][1]))
		} else {
			err = service.saveScraped(event)
			if err != nil {
				errs = append(errs, errors.Wrapf(err, "error importing Umweltkalender event: %s %s", event.Source, event.SourceId))
			}
		}
	}

	return errs
}

func (service *ScrapersService) Run() error {

	if !service.organizerId.Valid() {
		return errors.New("no valid organizerId set, exiting")
	}

	errs := service.scrapeUmweltKalender()
	if len(errs) > 0 {
		msg := ""
		for _, err := range errs {
			msg += "\n" + err.Error()
		}
		return errors.Errorf("error scraping Umweltkalender: %s", msg)
	}

	return nil
}
