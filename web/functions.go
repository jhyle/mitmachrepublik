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
		if t.Hour() == 0 && t.Minute() == 0 {
			return fmt.Sprintf("%02d.%02d.%04d", t.Day(), int(t.Month()), t.Year())
		} else {
			return fmt.Sprintf("%02d.%02d.%04d %02d:%02d", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute())
		}
	}
}

func longDatetimeFormat(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		if t.Hour() == 0 && t.Minute() == 0 {
			return fmt.Sprintf("%s, %02d.%02d.%04d", weekday[int(t.Weekday())], t.Day(), int(t.Month()), t.Year())
		} else {
			return fmt.Sprintf("%s, %02d.%02d.%04d %02d:%02d", weekday[int(t.Weekday())], t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute())
		}
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
	allowedTags       = []string{"h1", "h2", "h3", "h4", "h5", "h6", "div", "span", "hr", "p", "br", "b", "i", "strong", "em", "ol", "ul", "li", "a", "img", "table", "tbody", "tr", "td"}
	allowedAttributes = []string{"id", "class", "src", "href", "target", "title", "alt", "name", "rel", "style", "data-filename"}
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
			if !strings.HasSuffix(clipped, ".") {
				lastWhitespace := strings.LastIndexAny(clipped, " ,\t\r\n")
				if lastWhitespace > 0 {
					clipped = clipped[:lastWhitespace] + ".."
				} else {
					clipped += ".."
				}
			}
			break
		}
		runes++
	}

	return clipped
}

func strConcat(s []string) string {

	if len(s) == 0 {
		return ""
	} else if len(s) == 1 {
		return s[0]
	}
	
	return strings.Join(s[:len(s)-1], ", ") + " und " + s[len(s)-1]
}

func dates2RSSItems(dates []*Date) []rssItem {

	items := make([]rssItem, len(dates))
	for i, date := range dates {
		items[i] = rssItem{date.Id.Hex(), date.Title, citypartName(date.Addr) + ", " + datetimeFormat(date.Start) + " - " + date.PlainDescription(), date.Url()}
	}
	return items
}

func events2RSSItems(events []*Event) []rssItem {

	items := make([]rssItem, len(events))
	for i, event := range events {
		items[i] = rssItem{event.Id.Hex(), event.Title, citypartName(event.Addr) + ", " + datetimeFormat(event.Start) + " - " + event.PlainDescription(), event.Url()}
	}
	return items
}

func categoryIcon(categoryId int) string {

	return CategoryIconMap[categoryId]
}

func targetTitle(targetId int) string {

	return TargetIdMap[targetId]
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

	citypart := pcode2citypart[addr.Pcode]
	if isEmpty(citypart) {
		citypart = addr.City
	}
	return citypart
}

func encodePath(path string) string {

	return (&url.URL{Path: path}).String()
}

func eventSearchUrl(place string, targetIds, categoryIds, dateIds []int, radius int) string {

	targetNames := make([]string, len(targetIds))
	for i, id := range targetIds {
		targetNames[i] = TargetIdMap[id]
	}

	categoryNames := make([]string, len(categoryIds))
	for i, id := range categoryIds {
		categoryNames[i] = CategoryIdMap[id]
	}

	return place + "/" + strings.Join(int2Str(dateIds), ",") + "/" + strings.Join(int2Str(targetIds), ",")  + "/" + strings.Join(int2Str(categoryIds), ",") + "/" + strconv.Itoa(radius) + "/" + strings.Join(targetNames, ",") + "/" + strings.Join(categoryNames, ",") + "/0"
}

func eventSearchUrlWithQuery(place string, targetIds, categoryIds []int, dateIds []int, radius int, query string) string {

	return eventSearchUrl(place, targetIds, categoryIds, dateIds, radius) + "?query=" + query;
}

func simpleEventSearchUrl(place string) string {

	return eventSearchUrl(place, []int{0}, []int{0}, []int{0}, 0)
}

func categorySearchUrl(category int, place string) string {

	return eventSearchUrl(place, []int{0}, []int{category}, []int{0}, 0)
}

func targetSearchUrl(target int, place string) string {

	return eventSearchUrl(place, []int{target}, []int{0}, []int{0}, 0)
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

func inArray(a []int, n int) bool {

	for _, i := range a {
		if i == n {
			return true
		}
	}
	
	return false
}

func min(m, n int) int {

	if m < n {
		return m
	} else {
		return n
	}
}

func dayStart(now time.Time) time.Time {
	
	return time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local)
}

func dayEnd(now time.Time) time.Time {
	
	return time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, time.Local)
}

func timeSpans(dateIds []int) [][]time.Time {

	timeSpans := make([][]time.Time, len(dateIds))

	for i, dateId := range dateIds {
		now := time.Now()
		timespan := make([]time.Time, 2)

		if dateId == FromNow {
			timespan[0] = now
			timespan[1] = now
		} else if dateId == Today {
			timespan[0] = now
			timespan[1] = dayEnd(now)
		} else if dateId == Tomorrow {
			now = now.AddDate(0, 0, 1)
			timespan[0] = dayStart(now)
			timespan[1] = dayEnd(now)
		} else if dateId == AfterTomorrow {
			now = now.AddDate(0, 0, 1)
			now = now.AddDate(0, 0, 1)
			timespan[0] = dayStart(now)
			timespan[1] = dayEnd(now)
		} else if dateId == ThisWeek {
			timespan[0] = now
			for now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			timespan[1] = dayEnd(now)
		} else if dateId == NextWeekend {
			if now.Weekday() == time.Saturday || now.Weekday() == time.Sunday {
				timespan[0] = now
			} else {
				for now.Weekday() != time.Saturday && now.Weekday() != time.Sunday {
					now = now.AddDate(0, 0, 1)
				}
				timespan[0] = dayStart(now)
			}
			for now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			timespan[1] = dayEnd(now)
		} else if dateId == NextWeek {
			for now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			now = now.AddDate(0, 0, 1)
			timespan[0] = dayStart(now)
			for now.Weekday() != time.Sunday {
				now = now.AddDate(0, 0, 1)
			}
			timespan[1] = dayEnd(now)
		} else if dateId == TwoWeeks {
			timespan[0] = now
			for n := 0; n < 14; n++ {
				now = now.AddDate(0, 0, 1)
			}
			timespan[1] = dayEnd(now)
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

func intSlice(args ...int) []int {

    return args
}