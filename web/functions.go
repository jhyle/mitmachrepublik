package mmr

import (
	"fmt"
	"github.com/kennygrant/sanitize"
	"html/template"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func inc(i int) int {
	return i + 1
}

func dec(i int) int {
	return i - 1
}

func cut(s string, field int) string {

	return strings.Split(s, " ")[field]
}

var weekday map[int]string = map[int]string{
	int(time.Monday):    "Montag",
	int(time.Tuesday):   "Dienstag",
	int(time.Wednesday): "Mittwoch",
	int(time.Thursday):  "Donnerstag",
	int(time.Friday):    "Freitag",
	int(time.Saturday):  "Samstag",
	int(time.Sunday):    "Sonntag",
}

func dateFormat(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		return fmt.Sprintf("%s, %02d.%02d.%04d", weekday[int(t.Weekday())], t.Day(), int(t.Month()), t.Year())
	}
}

func timeFormat(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		return fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
	}
}

func datetimeFormat(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		return fmt.Sprintf("%02d.%02d.%04d %02d:%02d", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute())
	}
}

func iso8601Format(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		return fmt.Sprintf("%04d-%02d-%02dT%02d:%02d", t.Year(), int(t.Month()), t.Day(), t.Hour(), t.Minute())
	}
}

func noescape(s string) template.HTML {

	return template.HTML(s)
}

var (
	allowedTags = []string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "a", "img", "table", "tbody", "tr", "td"}
	allowedAttributes = []string{"id", "class", "src", "href", "title", "alt", "name", "rel", "style", "data-filename"}
)

func sanitizeHtml(s string) string {

	html, _ := sanitize.HTMLAllowing(strings.Replace(s, "data:", "data#", -1), allowedTags, allowedAttributes)
	return strings.Replace(html, "data#", "data:", -1)
}

func strClip(s string, n int) string {

	runes := 0
	clipped := s

	for index, _ := range s {
		if runes == n {
			clipped = s[:index]
			if strings.LastIndexAny(clipped, ".") != index {
				clipped = clipped[:strings.LastIndexAny(clipped, " ,\t\r\n")] + "..."
			}
			break
		}
		runes++
	}

	return clipped
}

func categoryIcon(categoryId int) string {

	return CategoryIconMap[categoryId]
}

func categoryTitle(categoryId int) string {

	return CategoryIdMap[categoryId]
}

func districtName(addr Address) string {

	district := pcode2district[addr.Pcode]
	if isEmpty(district) {
		district = addr.City
	}
	return district
}

func citypartName(addr Address) string {

	citypart := cpart2district[addr.Pcode]
	if isEmpty(citypart) {
		citypart = addr.City
	}
	return citypart
}

func encodePath(path string) string {

	return (&url.URL{Path: path}).String()
}

func eventSearchUrl(place string, categoryIds []int, dateIds []int, radius int) string {

	dateNames := make([]string, len(dateIds))
	for i, id := range dateIds {
		dateNames[i] = DateIdMap[id]
	}

	categoryNames := make([]string, len(categoryIds))
	for i, id := range categoryIds {
		categoryNames[i] = CategoryIdMap[id]
	}

	return place + "/" + strings.Join(dateNames, ",") + "/" + strings.Join(int2Str(categoryIds), ",") + "/" + strconv.Itoa(radius) + "/" + strings.Join(categoryNames, ",") + "/0"
}

func simpleEventSearchUrl(place string) string {

	return eventSearchUrl(place, []int{0}, []int{0}, 0)
}

func categorySearchUrl(category int, place string) string {

	return eventSearchUrl(place, []int{category}, []int{0}, 0)
}

func organizerSearchUrl(place string, categoryIds []int) string {

	categoryNames := make([]string, len(categoryIds))
	for i, id := range categoryIds {
		categoryNames[i] = CategoryIdMap[id]
	}
	return place + "/" + strings.Join(int2Str(categoryIds), ",") + "/" + strings.Join(categoryNames, ",") + "/0"
}

func simpleOrganizerSearchUrl(place string) string {

	return organizerSearchUrl(place, []int{0})
}

func str2Int(s []string) []int {

	a := make([]int, 0, len(s))

	for _, token := range s {
		n, err := strconv.Atoi(token)
		if err == nil {
			a = append(a, n)
		}
	}

	return a
}

func int2Str(i []int) []string {

	a := make([]string, len(i))

	for j, n := range i {
		a[j] = strconv.Itoa(n)
	}

	return a
}

func min(m, n int) int {

	if m < n {
		return m
	} else {
		return n
	}
}

func timeSpans(dateNames []string) [][]time.Time {

	timeSpans := make([][]time.Time, len(dateNames))

	for i, date := range dateNames {
		now := time.Now()
		timespan := make([]time.Time, 2)

		if date == "aktuell" {
			timespan[0] = now
			timespan[1] = now
		} else if date == "heute" {
			timespan[0] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			timespan[1] = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		} else if date == "morgen" {
			now = now.AddDate(0, 0, 1)
			timespan[0] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			timespan[1] = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		} else if date == "wochenende" {
			for now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			timespan[0] = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
			for now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			timespan[1] = time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
		}
		timeSpans[i] = timespan
	}

	return timeSpans
}

func pageCount(results, pageSize int) int {

	if results > 0 {
		results = results - 1
	}
	return (results / pageSize) + 1
}

func isEmpty(s string) bool {

	return len(strings.TrimSpace(s)) == 0
}
