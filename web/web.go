package mmr

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/pilu/traffic"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"strconv"
	"strings"
)

type (
	MmrApp struct {
		host         string
		port         int
		tpls         *Templates
		imgServer    string
		database     Database
		users        *UserService
		events       *EventService
		ga_code      string
		hostname     string
		emailAccount *EmailAccount
		locations    *LocationTree
		services     []Service
	}

	metaTags struct {
		Title    string
		FB_Title string
		FB_Image string
		FB_Descr string
		RSS      bool
	}

	emailAndPwd struct {
		Email string
		Pwd   string
	}

	sendMail struct {
		Name    string
		Email   string
		Subject string
		Text    string
	}

	appResult struct {
		Status int
		Error  error
		XML    string
		JSON   interface{}
	}

	rssItem struct {
		Id          string
		Title       string
		Description string
		Link        string
	}
)

const (
	register_subject = "Deine Registrierung bei mitmachrepublik.de"
	register_message = "Liebe/r Organisator/in von %s,\r\n\r\nvielen Dank für die Registrierung bei der Mitmach-Republik. Bitte bestätige Deine Registrierung, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://%s/approve/%s\r\n\r\nDas Team der Mitmach-Republik"
	password_subject = "Deine neue E-Mail-Adresse bei mitmachrepublik.de"
	password_message = "Liebe/r Organisator/in von %s,\r\n\r\nbitte bestätige Deine neue E-Mail-Adresse, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://%s/approve/%s\r\n\r\nDas Team der Mitmach-Republik"
	ga_dev           = "UA-61290824-1"
	ga_test          = "UA-61290824-2"
	ga_www           = "UA-61290824-3"
	google_api_key   = "AIzaSyAFzwmkGATzuHpcqV3g0yQEO77Vk66zXUM"
)

var (
	resultOK           = &appResult{Status: http.StatusOK}
	resultCreated      = &appResult{Status: http.StatusCreated}
	resultUnauthorized = &appResult{Status: http.StatusUnauthorized}
	resultBadRequest   = &appResult{Status: http.StatusBadRequest}
	resultNotFound     = &appResult{Status: http.StatusNotFound}
	resultConflict     = &appResult{Status: http.StatusConflict}
)

func NewMmrApp(env string, host string, port int, tplDir, imgServer, mongoUrl, dbName string) (*MmrApp, error) {

	database, err := NewMongoDb(mongoUrl, dbName)
	if err != nil {
		return nil, errors.New("init of MongoDB failed: " + err.Error())
	}

	users, err := NewUserService(database, "user")
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	events, err := NewEventService(database, "event", "date")
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	funcs := map[string]interface{}{
		"inc":                inc,
		"dec":                dec,
		"cut":                cut,
		"dateFormat":         dateFormat,
		"timeFormat":         timeFormat,
		"datetimeFormat":     datetimeFormat,
		"iso8601Format":      iso8601Format,
		"noescape":           noescape,
		"strClip":            strClip,
		"categoryIcon":       categoryIcon,
		"categoryTitle":      categoryTitle,
		"categorySearchUrl":  categorySearchUrl,
		"districtName":       districtName,
		"citypartName":       citypartName,
		"encodePath":         encodePath,
		"eventSearchUrl":     simpleEventSearchUrl,
		"organizerSearchUrl": simpleOrganizerSearchUrl,
	}

	tpls, err := NewTemplates(tplDir+string(os.PathSeparator)+"*.tpl", funcs)
	if err != nil {
		return nil, errors.New("init of templates failed: " + err.Error())
	}

	emailAccount := &EmailAccount{"smtp.gmail.com", 465, "mitmachrepublik", "mitmachen", &EmailAddress{"Mitmach-Republik", "mitmachrepublik@gmail.com"}}

	ga_code := ga_dev
	hostname := "dev.mitmachrepublik.de"
	if env == "www" {
		ga_code = ga_www
		hostname = "www.mitmachrepublik.de"
	} else if env == "test" {
		ga_code = ga_test
		hostname = "test.mitmachrepublik.de"
	}

	cities, err := events.Cities()
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	services := make([]Service, 0, 3)
	services = append(services, NewSessionService(60, database))
	services = append(services, NewUpdateRecurrencesService(3600, events))
	services = append(services, NewUnusedImgService(3600, database, imgServer))
	if env == "dev" {
		services = append(services, NewSpawnEventsService(3600, database, events, imgServer))
	}

	return &MmrApp{host, port, tpls, imgServer, database, users, events, ga_code, hostname, emailAccount, NewLocationTree(cities), services}, nil
}

func (app *MmrApp) view(tpl string, w traffic.ResponseWriter, meta *metaTags, data bson.M) *appResult {

	if data == nil {
		data = bson.M{}
	}
	data["meta"] = meta
	data["ga_code"] = app.ga_code
	data["hostname"] = app.hostname
	data["districts"] = DistrictMap
	data["categories"] = CategoryOrder
	data["categoryMap"] = CategoryMap
	data["googleApiKey"] = google_api_key

	err := app.tpls.Execute(tpl, w, data)
	if err != nil {
		return &appResult{Status: http.StatusInternalServerError, Error: err}
	}

	return resultOK
}

func (app *MmrApp) output(tpl string, w traffic.ResponseWriter, contentType string, meta *metaTags, data bson.M) *appResult {

	if data == nil {
		data = bson.M{}
	}
	data["meta"] = meta
	data["hostname"] = app.hostname

	w.Header().Set("Content-Type", contentType)
	err := app.tpls.Execute(tpl, w, data)
	if err != nil {
		return &appResult{Status: http.StatusInternalServerError, Error: err}
	}

	return resultOK
}

func (app *MmrApp) handle(w traffic.ResponseWriter, result *appResult) {

	if result.Error != nil {
		traffic.Logger().Print(result.Error.Error())
		if app.ga_code == ga_www {
			app.sendEmail(app.emailAccount.From, nil, "Fehlermeldung", result.Error.Error())
		}
	}

	if !w.Written() {
		if result == resultUnauthorized {
			w.Header().Set("Location", "/#login")
			w.WriteHeader(http.StatusFound)
		} else {
			w.WriteHeader(result.Status)
			if result == resultNotFound {
				app.staticPage(w, "notfound.tpl", "Seite nicht gefunden")
			}
		}
	}

	if result.JSON != nil && !w.BodyWritten() {
		w.WriteJSON(result.JSON)
	}
}

func (app *MmrApp) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	eventsPerRow := 4
	numberOfRows := 4
	pageSize := eventsPerRow * numberOfRows
	place := ""
	dateIds := []int{FromNow}

	meta := metaTags{
		"Willkommen in der Mitmach-Republik!",
		"Gemeinsam aktiv werden.",
		"http://" + app.hostname + "/images/mitmachrepublik.png",
		"Gemeinsam aktiv werden - hier findest Du Veranstaltungen und Organisationen zum Mitmachen. Finde Nachbarschaftstreffen, Vereine, gemeinnützige Initiativen und Ehrenämter in Deiner Umgebung. Mach mit bei gemeinsamen Projekten und Ideen!",
		true,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count(place, timeSpans(dateIds), nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		page := 0
		moreEvents := true
		events := make([]*Date, 0, eventsPerRow*2)
		for len(events) < eventsPerRow*numberOfRows && moreEvents {
			result, err := app.events.SearchDates(place, timeSpans(dateIds), nil, true, page, pageSize, "start")
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}
			if result.GetSize() == 0 {
				moreEvents = false
			} else {
				for _, date := range result.Dates {
					alreadyIncluded := false
					for _, event := range events {
						if event.EventId == date.EventId {
							alreadyIncluded = true
							break
						}
					}
					if !alreadyIncluded && len(events) < eventsPerRow*numberOfRows {
						events = append(events, date)
					}
				}
			}
			page++
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": dates2RSSItems(events)})
		} else {
			var dates [][]*Date
			if len(events) > 0 {
				dates = make([][]*Date, ((len(events)-1)/eventsPerRow)+1)
				for i := range dates {

					rowSize := len(events) - i*eventsPerRow
					if rowSize > eventsPerRow {
						rowSize = eventsPerRow
					}
					dates[i] = make([]*Date, rowSize)

					for j := 0; j < rowSize; j++ {
						dates[i][j] = events[i*eventsPerRow+j]
					}
				}
			} else {
				dates = make([][]*Date, 0)
			}
			return app.view("start.tpl", w, &meta, bson.M{"events": dates, "eventCnt": eventCnt, "organizerCnt": organizerCnt, "dates": DateOrder, "dateMap": DateIdMap})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) approvePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{
		"Registrierung bestätigen - Mitmach-Republik",
		"Registrierung bestätigen",
		"http://" + app.hostname + "/images/mitmachrepublik.png",
		"",
		false,
	}

	result := func() *appResult {

		var err error = nil

		if bson.IsObjectIdHex(r.Param("id")) {
			var user *User
			user, err = app.users.Load(bson.ObjectIdHex(r.Param("id")))
			if err == nil {
				user.Approved = true
				err := app.users.Store(user)
				if err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}

				events, err := app.events.FindEventsOfUser(user.Id, "start")
				if err == nil {
					for _, event := range events {
						err = app.events.SyncDates(&event)
						if err != nil {
							break
						}
					}
				}
				if err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}
			}
		} else {
			err = errors.New("No organizer id given.")
		}

		return app.view("approve.tpl", w, &meta, bson.M{"approved": err == nil})
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 5
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	if r.Param("fmt") == "RSS" {
		pageSize = 1000
		page = 0
	}

	place := r.Param("place")
	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	var dateIds []int
	if r.Param("dates") == "aktuell" {
		dateIds = []int{FromNow}
	} else {
		dateIds = str2Int(strings.Split(r.Param("dates"), ","))
	}
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	title := "Aktuelle Veranstaltungen"
	if len(categoryIds) > 0 {
		categories := make([]string, len(categoryIds))
		for i, categoryId := range categoryIds {
			categories[i] = CategoryIdMap[categoryId]
		}
		title += " aus " + strConcat(categories)
	}
	if !isEmpty(place) {
		title += " in " + place
	}

	meta := metaTags{
		title + " - Mitmach-Republik",
		title,
		"http://" + app.hostname + "/images/mitmachrepublik.png",
		"Veranstaltungen zum Mitmachen! Heute, morgen oder am nächsten Wochenende - finde Veranstaltungen in den Kategorien " + strConcat(CategoryOrder) + ".",
		true,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count(place, timeSpans(dateIds), categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.events.SearchDates(place, timeSpans(dateIds), categoryIds, false, page, pageSize, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerNames := make(map[bson.ObjectId]string)
		for _, date := range result.Dates {
			if _, found := organizerNames[date.OrganizerId]; !found {
				user, err := app.users.Load(date.OrganizerId)
				if err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}
				organizerNames[date.OrganizerId] = user.Name
			}
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": dates2RSSItems(result.Dates)})
		} else {
			pageCount := pageCount(result.Count, pageSize)
			pages := make([]int, pageCount)
			for i := 0; i < pageCount; i++ {
				pages[i] = i
			}
			maxPage := pageCount - 1
			return app.view("events.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Dates, "organizerNames": organizerNames, "place": place, "radius": radius, "dates": DateOrder, "dateMap": DateIdMap, "dateIds": dateIds, "categoryIds": categoryIds})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventPage(w traffic.ResponseWriter, r *traffic.Request) {

	radius := 2
	dateIds := []int{FromNow}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		date, err := app.events.LoadDate(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		recurrences, err := app.events.FindDatesOfEvent(date.EventId, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizer, err := app.users.Load(date.OrganizerId)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		place := citypartName(date.Addr)

		eventCnt, err := app.events.Count(place, timeSpans(dateIds), nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		imageUrl := ""
		if !isEmpty(date.Image) {
			imageUrl = "http://" + app.hostname + "/bild/" + date.Image
		}

		title := "Veranstaltung " + date.Title
		title += " am " + dateFormat(date.Start)
		if !isEmpty(place) {
			title += " in " + place
		}
		if !isEmpty(date.Addr.Name) {
			title += " (" + date.Addr.Name + ")"
		}

		meta := metaTags{
			title + " - Mitmach-Republik",
			title,
			imageUrl,
			strClip(date.PlainDescription(), 160),
			false,
		}

		return app.view("event.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "place": place, "radius": radius, "event": date, "organizer": organizer, "recurrences": recurrences})
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendEventPage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"Veranstaltung empfehlen - Mitmach-Republik", "", "", "", false}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		event, err := app.events.LoadDate(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		return app.view("sendevent.tpl", w, &meta, bson.M{"event": &event})
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizersPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 5
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	place := r.Param("place")
	dateIds := []int{FromNow}
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	title := "Gemeinschaftliche Organisatoren"
	if len(categoryIds) > 0 {
		categories := make([]string, len(categoryIds))
		for i, categoryId := range categoryIds {
			categories[i] = CategoryIdMap[categoryId]
		}
		title += " aus " + strConcat(categories)
	}
	if !isEmpty(place) {
		title += " in " + place
	}

	meta := metaTags{
		title + " - Mitmach-Republik",
		title,
		"http://" + app.hostname + "/images/mitmachrepublik.png",
		"Mitmacher gesucht! Finde Organisatoren in den Kategorien " + strConcat(CategoryOrder) + ".",
		false,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count(place, timeSpans(dateIds), categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.users.Search(place, categoryIds, page, pageSize, "name")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		pageCount := pageCount(result.Count, pageSize)
		pages := make([]int, pageCount)
		for i := 0; i < pageCount; i++ {
			pages[i] = i
		}
		maxPage := pageCount - 1

		return app.view("organizers.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "organizers": result.Organizers, "place": place, "radius": radius, "categoryIds": categoryIds})
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizerPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 5
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	if r.Param("fmt") == "RSS" {
		pageSize = 1000
		page = 0
	}

	radius := 2
	dateIds := []int{FromNow}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		organizer, err := app.users.Load(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		place := citypartName(organizer.Addr)

		eventCnt, err := app.events.Count(place, timeSpans(dateIds), nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.events.SearchDatesOfUser(organizer.Id, page, pageSize, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerNames := map[bson.ObjectId]string{organizer.Id: organizer.Name}

		imageUrl := ""
		if !isEmpty(organizer.Image) {
			imageUrl = "http://" + app.hostname + "/bild/" + organizer.Image
		}
		title := "Gemeinschaftliche Veranstaltungen"
		if organizer.Name != "Mitmach-Republik" {
			title += " von " + organizer.Name
			if !isEmpty(place) {
				title += " aus " + place
			}
		}
		meta := metaTags{
			title + " - Mitmach-Republik",
			title,
			imageUrl,
			strClip(organizer.PlainDescription(), 160),
			true,
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": dates2RSSItems(result.Dates)})
		} else {
			pageCount := pageCount(result.Count, pageSize)
			pages := make([]int, pageCount)
			for i := 0; i < pageCount; i++ {
				pages[i] = i
			}
			maxPage := pageCount - 1
			return app.view("organizer.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Dates, "organizerNames": organizerNames, "place": place, "radius": radius, "organizer": organizer})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) checkSession(r *Request) (*User, error) {

	sessionId, err := r.ReadSessionId()
	if err != nil {
		return nil, err
	}

	user, err := app.database.LoadUserBySessionId(sessionId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (app *MmrApp) adminPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 5
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	query := r.Param("q")
	meta := metaTags{"Verwaltung - Mitmach-Republik", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		result, err := app.events.SearchEventsOfUser(user.Id, query, page, pageSize, "-start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		pageCount := result.Count / pageSize
		if pageCount == 0 {
			pageCount = 1
		}
		pages := make([]int, pageCount)
		for i := 0; i < pageCount; i++ {
			pages[i] = i
		}
		maxPage := pageCount - 1

		return app.view("admin.tpl", w, &meta, bson.M{"user": user, "query": query, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events})
	}()

	app.handle(w, result)
}

func (app *MmrApp) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"Profil bearbeiten - Mitmach-Republik", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		return app.view("profile.tpl", w, &meta, bson.M{"user": user})
	}()

	app.handle(w, result)
}

func (app *MmrApp) passwordPage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"E-Mail-Adresse oder Kennwort ändern - Mitmach-Republik", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		return app.view("password.tpl", w, &meta, bson.M{"user": user})
	}()

	app.handle(w, result)
}

func (app *MmrApp) editEventPage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"Veranstaltung bearbeiten - Mitmach-Republik", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		var event *Event
		if bson.IsObjectIdHex(r.Param("id")) {

			event, err = app.events.Load(bson.ObjectIdHex(r.Param("id")))
			if err != nil {
				return &appResult{Status: http.StatusNotFound, Error: err}
			}

			if event.OrganizerId != user.Id {
				return resultBadRequest
			}
		} else if bson.IsObjectIdHex(r.Param("copy")) {

			oldEvent, err := app.events.Load(bson.ObjectIdHex(r.Param("copy")))
			if err != nil {
				return &appResult{Status: http.StatusNotFound, Error: err}
			}

			event = new(Event)
			event.OrganizerId = user.Id
			event.Title = oldEvent.Title
			event.Start = oldEvent.Start
			event.End = oldEvent.End
			event.Rsvp = oldEvent.Rsvp
			event.Image = oldEvent.Image
			event.Categories = oldEvent.Categories
			event.Descr = oldEvent.Descr
			event.Addr = oldEvent.Addr
			event.Web = oldEvent.Web
			event.Recurrency = oldEvent.Recurrency
			event.Monthly = oldEvent.Monthly
			event.Weekly = oldEvent.Weekly
		} else {
			event = new(Event)
		}

		return app.view("event_edit.tpl", w, &meta, bson.M{"user": user, "event": event})
	}()

	app.handle(w, result)
}

func (app *MmrApp) staticPage(w traffic.ResponseWriter, template, headline string) {

	meta := metaTags{headline + " - Mitmach-Republik", headline, "", "", false}

	result := func() *appResult {
		return app.view(template, w, &meta, nil)
	}()

	app.handle(w, result)
}

func (app *MmrApp) searchHandler(w traffic.ResponseWriter, r *traffic.Request) {

	place := strings.Trim(r.PostFormValue("place"), " ")

	radius, err := strconv.Atoi(r.PostFormValue("radius"))
	if err != nil {
		radius = 0
	}

	categoryIds := str2Int(r.Form["category"])
	if len(categoryIds) == 0 {
		categoryIds = append(categoryIds, 0)
	}

	dateIds := str2Int(r.Form["date"])
	if len(dateIds) == 0 {
		dateIds = append(dateIds, 0)
	}

	path := r.PostFormValue("search")
	if path == "organizers" {
		path = "/veranstalter/" + organizerSearchUrl(place, categoryIds)
	} else {
		path = "/veranstaltungen/" + eventSearchUrl(place, categoryIds, dateIds, radius)
	}

	w.Header().Set("Location", path)
	w.WriteHeader(http.StatusFound)
}

func (app *MmrApp) uploadHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		file, _, err := r.FormFile("file")
		if err != nil {
			return resultBadRequest
		}

		filename := uuid.New() + ".jpg"
		resp, err := http.Post(app.imgServer+"/"+filename, "image/jpeg", file)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		if resp.StatusCode == http.StatusOK {
			return &appResult{Status: http.StatusOK, JSON: filename}
		} else {
			return &appResult{Status: resp.StatusCode}
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendEmail(to, replyTo *EmailAddress, subject, body string) error {

	return SendEmail(app.emailAccount, to, replyTo, subject, body)
}

func (app *MmrApp) registerHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		user, err := (&Request{r}).ReadUser()
		if err != nil || !user.AGBs {
			return resultBadRequest
		}

		err = app.users.Validate(user)
		if err != nil {
			return resultConflict
		}

		user.Id = bson.NewObjectId()
		err = app.users.Store(user)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = app.sendEmail(&EmailAddress{user.Name, user.Email}, nil, register_subject, fmt.Sprintf(register_message, user.Name, app.hostname, user.Id.Hex()))
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		id, err := app.database.CreateSession(user.Id)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusCreated, JSON: id}
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendCheckMailHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		request := &Request{r}

		user, err := app.checkSession(request)
		if err != nil {
			return resultUnauthorized
		}

		err = app.sendEmail(&EmailAddress{user.Name, user.Email}, nil, register_subject, fmt.Sprintf(register_message, user.Name, app.hostname, user.Id.Hex()))
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) profileHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		request := &Request{r}

		user, err := app.checkSession(request)
		if err != nil {
			return resultUnauthorized
		}

		data, err := request.ReadUser()
		if err != nil {
			return resultBadRequest
		}

		user.Name = data.Name
		user.Image = data.Image
		user.Categories = data.Categories
		user.Descr = data.Descr
		user.Web = data.Web
		user.Addr.Street = data.Addr.Street
		user.Addr.Pcode = data.Addr.Pcode
		user.Addr.City = data.Addr.City

		err = app.users.Store(user)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) passwordHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		request := &Request{r}

		user, err := app.checkSession(request)
		if err != nil {
			return resultUnauthorized
		}

		data, err := request.ReadEmailAndPwd()
		if err != nil {
			return resultBadRequest
		}

		if !isEmpty(data.Pwd) {
			user.Pwd = data.Pwd
		}

		if !isEmpty(data.Email) && data.Email != user.Email {

			user.Email = data.Email
			user.Approved = false
			err := app.events.DeleteDatesOfUser(user.Id)
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}

			err = app.sendEmail(&EmailAddress{user.Name, user.Email}, nil, password_subject, fmt.Sprintf(password_message, user.Name, app.hostname, user.Id.Hex()))
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}
		}

		err = app.users.Store(user)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) unregisterHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		user, err := app.checkSession(&Request{r})
		if err != nil {
			return resultUnauthorized
		}

		err = app.events.DeleteOfUser(user.Id)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = app.users.Delete(user.Id)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
	app.logoutHandler(w, r)
}

func (app *MmrApp) eventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		request := &Request{r}

		user, err := app.checkSession(request)
		if err != nil {
			return resultUnauthorized
		}

		var data Event
		event := &data
		err = request.ReadJson(event)
		if err != nil {
			return resultBadRequest
		}

		created := false
		if event.Id.Valid() {
			oldEvent, err := app.events.Load(event.Id)
			if err != nil || oldEvent.OrganizerId != user.Id {
				return resultBadRequest
			}
		} else {
			created = true
			event.Id = bson.NewObjectId()
		}
		event.OrganizerId = user.Id

		err = app.events.Store(event, user.Approved)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		app.locations.Add(event.Addr.City)

		if created {
			return resultCreated
		} else {
			return resultOK
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendEventMailHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		form, err := (&Request{r}).ReadSendMail()
		if err != nil {
			return resultBadRequest
		}

		err = app.sendEmail(&EmailAddress{form.Name, form.Email}, nil, form.Subject, form.Text)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendContactMailHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		form, err := (&Request{r}).ReadSendMail()
		if err != nil {
			return resultBadRequest
		}

		err = app.sendEmail(app.emailAccount.From, &EmailAddress{form.Name, form.Email}, form.Subject, form.Text)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) locationHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		locations := app.locations.Autocomplete(r.Param("location"))
		return &appResult{Status: http.StatusOK, JSON: locations}
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		dateIds := str2Int(strings.Split(r.Param("dateIds"), ","))
		categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

		cnt, err := app.events.Count(r.Param("place"), timeSpans(dateIds), categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusOK, JSON: cnt}
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizerCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

		cnt, err := app.users.Count(r.Param("place"), categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusOK, JSON: cnt}
	}()

	app.handle(w, result)
}

func (app *MmrApp) loginHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		form, err := (&Request{r}).ReadEmailAndPwd()
		if err != nil {
			return resultBadRequest
		}

		user, err := app.database.LoadUserByEmailAndPassword(form.Email, form.Pwd)
		if err != nil {
			return resultNotFound
		}

		id, err := app.database.CreateSession(user.Id)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusCreated, JSON: id}
	}()

	app.handle(w, result)
}

func (app *MmrApp) logoutHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		sessionId, err := (&Request{r}).ReadSessionId()
		if err != nil {
			return resultBadRequest
		}

		err = app.database.RemoveSession(sessionId)
		if err != nil {
			return resultNotFound
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) deleteEventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		user, err := app.checkSession(&Request{r})
		if err != nil {
			return resultUnauthorized
		}

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultBadRequest
		}

		event, err := app.events.Load(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultOK
		}

		if event.OrganizerId != user.Id {
			return resultBadRequest
		}

		err = app.events.Delete(event.Id)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		// TODO delete dates

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) Start() {

	for _, service := range app.services {
		service.Start()
	}

	traffic.SetHost(app.host)
	traffic.SetPort(app.port)
	router := traffic.New()

	router.Get("/", app.startPage)
	router.Get("/veranstalter/verwaltung/kennwort", app.passwordPage)
	router.Get("/veranstalter/verwaltung/profil", app.profilePage)
	router.Get("/veranstalter/verwaltung/veranstaltung", app.editEventPage)
	router.Get("/veranstalter/verwaltung/veranstaltung/:id", app.editEventPage)
	router.Get("/veranstalter/verwaltung/:page", app.adminPage)
	router.Get("/veranstaltungen/:place/:dates/:categoryIds/:radius/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen//:dates/:categoryIds/:radius/:categories/:page", app.eventsPage)
	router.Get("/veranstaltung/:categories/:date/:id/:title", app.eventPage)
	router.Get("/veranstaltung//:date/:id/:title", app.eventPage)
	router.Get("/veranstalter/:place/:categoryIds/:categories/:page", app.organizersPage)
	router.Get("/veranstalter//:categoryIds/:categories/:page", app.organizersPage)
	router.Get("/veranstalter/:id/:title/:page", app.organizerPage)

	router.Get("/impressum", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "impressum.tpl", "Impressum") })
	router.Get("/disclaimer", func(w traffic.ResponseWriter, r *traffic.Request) {
		app.staticPage(w, "disclaimer.tpl", "Haftungsausschluss (Disclaimer)")
	})
	router.Get("/datenschutz", func(w traffic.ResponseWriter, r *traffic.Request) {
		app.staticPage(w, "datenschutz.tpl", "Datenschutzerklärung")
	})
	router.Get("/agbs", func(w traffic.ResponseWriter, r *traffic.Request) {
		app.staticPage(w, "agbs.tpl", "Allgemeine Geschäftsbedingungen")
	})
	router.Get("/dialog/contact", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "contact.tpl", "") })
	router.Get("/dialog/registered", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "registered.tpl", "") })
	router.Get("/dialog/login", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "login.tpl", "") })
	router.Get("/dialog/register", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "register.tpl", "") })
	router.Get("/dialog/sendevent/:id", app.sendEventPage)

	router.Post("/suche", app.searchHandler)
	router.Post("/upload", app.uploadHandler)
	router.Post("/register", app.registerHandler)
	router.Post("/sendcheckmail", app.sendCheckMailHandler)
	router.Post("/password", app.passwordHandler)
	router.Post("/profile", app.profileHandler)
	router.Post("/unregister", app.unregisterHandler)
	router.Post("/event", app.eventHandler)
	router.Post("/sendeventmail", app.sendEventMailHandler)
	router.Post("/sendcontactmail", app.sendContactMailHandler)
	router.Get("/location/:location", app.locationHandler)
	router.Get("/eventcount/:place/:dateIds/:categoryIds", app.eventCountHandler)
	router.Get("/eventcount//:dateIds/:categoryIds", app.eventCountHandler)
	router.Get("/organizercount/:place/:categoryIds", app.organizerCountHandler)
	router.Get("/organizercount//:categoryIds", app.organizerCountHandler)
	router.Post("/login", app.loginHandler)
	router.Post("/logout", app.logoutHandler)

	router.Get("/approve/:id", app.approvePage)
	router.Delete("/event/:id", app.deleteEventHandler)

	router.Run()
}
