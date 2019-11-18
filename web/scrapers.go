package mmr

import (
	"bytes"
	"fmt"
	"html"
	"io/ioutil"
	"net/http"
	"net/url"
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
	UK_SOURCE = "umweltkalender"

	TAG_H1_REGEX        = "<h1[^>]*>(.*)</h1"
	POSTCODE_CITY_REGEX = "(\\d{5})\\s+(.*)"
	POSTCODE_REGEX      = "(\\d{5})"

	umweltKalenderBaseUrl     = "https://www.umweltkalender-berlin.de"
	umweltKalenderListUrl     = umweltKalenderBaseUrl + "/angebote/filter?ps=1"
	umweltKalenderListRequest = `{"arten":["MIMA"],"ort":[],"zielgr":[],"sehb":[],"roll":[],"zeitraum":"zz","plain":"","um_km":"","um_plz":"","ab_tag":"","ab_monat":"","ab_jahr":""}`
)

var (
	ANMELDUNG_ERFORDERLICH = []byte("ANMELDUNG ERFORDERLICH!")

	tag_h1        *regexp.Regexp = regexp.MustCompile(TAG_H1_REGEX)
	postcode_city *regexp.Regexp = regexp.MustCompile(POSTCODE_CITY_REGEX)
	postcode      *regexp.Regexp = regexp.MustCompile(POSTCODE_REGEX)

	umweltKalenderDetailsLink      = regexp.MustCompile(`(?ms)grid-item teaser.*?href="([^\"]+)`)
	umweltKalenderIdAndDateFromUrl = regexp.MustCompile(`/(\d+)\?dat=(\d{4})-(\d{2})-(\d{2})`)
	umweltKalenderStartAndEnd      = regexp.MustCompile(`(?ms)date_detail.*?</div>(\d{2}):(\d{2}) - (\d{2}):(\d{2})`)
	umweltKalenderDescription      = regexp.MustCompile(`(?ms)<div class="read-more-content">\s+(.*)</div>\s+<div class="accordeon_hl read-more-trigger">MEHR ANZEIGEN`)
	umweltKalenderDescription2     = regexp.MustCompile(`(?ms)<hr class="js-hide-termine" />\s+(.*)\s+<hr />`)
	umweltKalenderLocation         = regexp.MustCompile(`(?ms)<strong>Treffpunkt:</strong><br>\s+<div>([^<]+)</div>`)
	umweltKalenderTargets          = regexp.MustCompile(`(?ms)<strong>Für:</strong><br>\s+<div>([^<]+)</div>`)
	umweltKalenderLinks            = regexp.MustCompile(`<a href="(/angebote/details/\d+?dat=\d{4}-\d{2}-\d{2})"`)

	importedTags       = []string{"div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "table", "tbody", "tr", "td"}
	importedAttributes = []string{"title"}
)

func NewScrapersService(hour int, email *EmailAccount, events *EventService, organizerId bson.ObjectId) Service {

	return &ScrapersService{NewBasicService("NewSraperService", hour, email), events, organizerId}
}

func (service *ScrapersService) Start() {

	service.start(service.Run)
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

func (service *ScrapersService) loadUmweltKalenderListPage() ([]byte, error) {

	request := url.Values{"filterJson": {umweltKalenderListRequest}}

	response, err := http.PostForm(umweltKalenderListUrl, request)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load Umweltkalender list page")
	}
	defer response.Body.Close()

	page, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrap(err, "error reading Umweltkalender list page")
	}

	return page, nil

}

func (service *ScrapersService) getIdAndDateFromUmweltKalendarDetailsLink(path string) (bool, string, time.Time) {

	match := umweltKalenderIdAndDateFromUrl.FindStringSubmatch(path)
	if match == nil {
		return false, "", time.Time{}
	}

	id := string(match[1])
	year, _ := strconv.ParseUint(string(match[2]), 10, 32)
	month, _ := strconv.ParseUint(string(match[3]), 10, 32)
	day, _ := strconv.ParseUint(string(match[4]), 10, 32)
	date := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)

	return true, id, date
}

func (service *ScrapersService) getStartAndEndFromUmweltKalenderPage(page []byte) (bool, time.Duration, time.Duration) {

	match := umweltKalenderStartAndEnd.FindSubmatch(page)
	if match == nil {
		return false, 0, 0
	}

	startHour, _ := strconv.ParseUint(string(match[1]), 10, 32)
	startMinute, _ := strconv.ParseUint(string(match[2]), 10, 32)
	start := time.Duration(startHour)*time.Hour + time.Duration(startMinute)*time.Minute

	endHour, _ := strconv.ParseUint(string(match[3]), 10, 32)
	endMinute, _ := strconv.ParseUint(string(match[4]), 10, 32)
	end := time.Duration(endHour)*time.Hour + time.Duration(endMinute)*time.Minute

	return true, start, end
}

func (service *ScrapersService) getTitleAndDescriptionFromUmweltKalenderPage(page []byte) (string, string, error) {

	h1 := tag_h1.FindSubmatch(page)
	if h1 == nil {
		return "", "", errors.New("could not find title of Umweltkalender event")
	}
	title := html.UnescapeString(string(h1[1]))

	description := umweltKalenderDescription.FindSubmatch(page)
	if description == nil {
		return "", "", errors.New("could not find description of Umweltkalender event")
	}

	descrHTML, err := sanitize.HTMLAllowing(string(description[1]), importedTags, importedAttributes)
	if err != nil {
		return "", "", err
	}

	description2 := umweltKalenderDescription2.FindSubmatch(page)
	if description2 != nil {
		descrHTML2, err := sanitize.HTMLAllowing(string(description2[1]), importedTags, importedAttributes)
		if err != nil {
			return "", "", err
		}
		descrHTML += "<p>" + descrHTML2 + "</p>"
	}
	descrHTML += "<p class=\"small\">Text von der Veranstaltungswebseite übernommen</p>"

	return title, descrHTML, nil
}

func (service *ScrapersService) getAddressFromUmweltKalenderPage(page []byte) (Address, error) {

	var address Address

	location := umweltKalenderLocation.FindSubmatch(page)
	if location == nil {
		return address, errors.New("could not find address of Umweltkalender event")
	}
	line := html.UnescapeString(string(location[1]))

	postcode_part := -1
	parts := strings.Split(line, ", ")
	for i := range parts {
		if postcode.MatchString(parts[i]) {
			postcode_part = i
			break
		}
	}
	if postcode_part > 0 {
		address.Street = strings.TrimSpace(parts[postcode_part-1])
	}
	if postcode_part > -1 {
		district := postcode_city.FindStringSubmatch(parts[postcode_part])
		if district != nil {
			address.Pcode = district[1]
			address.City = strings.TrimSpace(district[2])
		} else {
			pcode := postcode.FindStringSubmatch(parts[postcode_part])
			if pcode == nil {
				return address, errors.New("could not find postcode or city")
			}
			address.Pcode = pcode[1]
			address.City = "Berlin"
		}
	}

	return address, nil
}

func (service *ScrapersService) getTargetsFromUmweltKalenderPage(page []byte) ([]int, error) {

	var targetIds []int

	targets := umweltKalenderTargets.FindSubmatch(page)
	if targets == nil {
		return nil, errors.New("could not find targets of Umweltkalender event")
	}
	line := html.UnescapeString(string(targets[1]))

	for _, target := range strings.Split(line, ", ") {
		if targetId, exists := TargetMap[strings.TrimSpace(target)]; exists {
			targetIds = append(targetIds, targetId)
		}
	}

	return targetIds, nil
}

func (service *ScrapersService) findRelatedUmweltKalenderEvents(page []byte) ([]string, error) {

	var related []string

	links := umweltKalenderLinks.FindAllSubmatch(page, -1)
	for _, link := range links {
		related = append(related, string(link[1]))
	}

	return related, nil
}

func (service *ScrapersService) loadUmweltKalenderDetailsPage(path string) ([]byte, error) {

	url := umweltKalenderBaseUrl + path

	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error loading url: %s", url)
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "error reading body of url: %s", url)
	}

	return body, nil
}

func (service *ScrapersService) importUmweltKalenderEvent(id string, date time.Time, link string) ([]string, error) {

	page, err := service.loadUmweltKalenderDetailsPage(link)
	if err != nil {
		return nil, err
	}

	hasTime, start, end := service.getStartAndEndFromUmweltKalenderPage(page)
	if !hasTime {
		return nil, nil
	}
	startDate := date.Add(start)
	endDate := date.Add(end)

	title, description, err := service.getTitleAndDescriptionFromUmweltKalenderPage(page)
	if err != nil {
		return nil, err
	}

	address, err := service.getAddressFromUmweltKalenderPage(page)
	if err != nil {
		return nil, err
	}

	targets, err := service.getTargetsFromUmweltKalenderPage(page)
	if err != nil {
		return nil, err
	}

	categories := []int{CategoryMap["Leute treffen"], CategoryMap["Umweltschutz"]}
	rsvp := bytes.Contains(page, ANMELDUNG_ERFORDERLICH)

	event := &Event{
		Source:     UK_SOURCE,
		SourceId:   fmt.Sprintf("%s#%d-%d-%d", id, date.Year(), date.Month(), date.Day()),
		Start:      startDate,
		End:        endDate,
		Title:      title,
		Descr:      description,
		Addr:       address,
		Web:        umweltKalenderBaseUrl + link,
		Targets:    targets,
		Categories: categories,
		Rsvp:       rsvp,
	}

	err = service.saveScraped(event)
	if err != nil {
		return nil, err
	}

	related, err := service.findRelatedUmweltKalenderEvents(page)
	if err != nil {
		return nil, err
	}

	return related, nil
}

func (service *ScrapersService) scrapeUmweltKalender2() []error {

	page, err := service.loadUmweltKalenderListPage()
	if err != nil {
		return []error{err}
	}

	detailsLinks := umweltKalenderDetailsLink.FindAllSubmatch(page, -1)
	if detailsLinks == nil {
		return []error{errors.New("no events found at Umweltkalender")}
	}

	var errs []error
	handled := map[string]bool{}

	for i := range detailsLinks {

		link := string(detailsLinks[i][1])
		if handled[link] {
			continue
		}
		handled[link] = true

		hasDate, id, date := service.getIdAndDateFromUmweltKalendarDetailsLink(link)
		if !hasDate {
			continue
		}

		moreLinks, err := service.importUmweltKalenderEvent(id, date, link)
		if err != nil {
			errs = append(errs, errors.Wrapf(err, "error importing %s", link))
		}

		for _, link := range moreLinks {

			if handled[link] {
				continue
			}
			handled[link] = true

			hasDate, id, date := service.getIdAndDateFromUmweltKalendarDetailsLink(link)
			if !hasDate {
				continue
			}

			_, err = service.importUmweltKalenderEvent(id, date, link)
			if err != nil {
				errs = append(errs, errors.Wrapf(err, "error importing %s", link))
			}
		}
	}

	return errs
}

func (service *ScrapersService) Run() error {

	if !service.organizerId.Valid() {
		return errors.New("no valid organizerId set, exiting")
	}

	errs := service.scrapeUmweltKalender2()
	if len(errs) > 0 {
		msg := ""
		for _, err := range errs {
			msg += "\n" + err.Error()
		}
		return errors.Errorf("error scraping Umweltkalender: %s", msg)
	}

	return nil
}
