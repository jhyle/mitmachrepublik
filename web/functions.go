package mmr

import (
	"errors"
	"fmt"
	"labix.org/v2/mgo/bson"
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
		return fmt.Sprintf("%02d.%02d Uhr", t.Hour(), t.Minute())
	}
}

func datetimeFormat(t time.Time) string {

	if t.IsZero() {
		return ""
	} else {
		return fmt.Sprintf("%02d.%02d.%04d %02d.%02d Uhr", t.Day(), int(t.Month()), t.Year(), t.Hour(), t.Minute())
	}
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

func eventUrl(event *Event) string {

	categoryNames := make([]string, len(event.Categories))
	for i, id := range event.Categories {
		categoryNames[i] = CategoryIdMap[id]
	}

	return strings.Join(categoryNames, ",") + "/" + dateFormat(event.Start) + "/" + event.Id.Hex() + "/" + event.Title
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

func organizerUrl(organizer User) string {

	return organizer.Id.Hex() + "/" + organizer.Addr.Name + "/0"
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

func buildQuery(place string, dates [][]time.Time, categoryIds []int) bson.M {

	query := make([]bson.M, 0, 3)

	if len(place) > 0 {
		postcodes := Postcodes(place)
		placesQuery := make([]bson.M, len(postcodes)+1)
		for i, postcode := range postcodes {
			placesQuery[i] = bson.M{"addr.pcode": postcode}
		}
		placesQuery[len(postcodes)] = bson.M{"addr.city": place}
		query = append(query, bson.M{"$or": placesQuery})
	}

	if len(dates) > 0 {
		datesQuery := make([]bson.M, len(dates))
		for i, timespan := range dates {
			rangeQuery := make(bson.M)
			rangeQuery["$gte"] = timespan[0]
			if timespan[1] != timespan[0] {
				rangeQuery["$lt"] = timespan[1]
			}
			datesQuery[i] = bson.M{"start": rangeQuery}
		}
		query = append(query, bson.M{"$or": datesQuery})
	}

	if len(categoryIds) > 0 && categoryIds[0] != 0 {
		categoriesQuery := make([]bson.M, len(categoryIds))
		for i, categoryId := range categoryIds {
			categoriesQuery[i] = bson.M{"categories": categoryId}
		}
		query = append(query, bson.M{"$or": categoriesQuery})
	}

	return bson.M{"$and": query}
}

func isEmpty(s string) bool {

	return len(strings.TrimSpace(s)) == 0
}

func validateUser(db Database, user *User) error {

	table := db.Table("user")

	var result []User
	err := table.Find(bson.M{"email": user.Email}, &result)
	if err != nil {
		return err
	}

	for i := 0; i < len(result); i++ {
		if result[i].Email == user.Email && result[i].Id != user.Id {
			return errors.New("Email address is already in use.")
		}
	}

	return nil
}
