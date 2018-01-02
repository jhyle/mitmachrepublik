package mmr

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/pilu/traffic"
	"gopkg.in/mgo.v2/bson"
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
		alerts       *AlertService
		ga_code      string
		hostname     string
		emailAccount *EmailAccount
		locations    *LocationTree
		services     []Service
		fbAppSecret  string
		fbUser       string
		fbPass       string
	}

	metaTags struct {
		Title    string
		Descr    string
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

	sendPassword struct {
		Email string
	}

	appResult struct {
		Status      int
		Error       error
		RedirectUrl string
		XML         string
		JSON        interface{}
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
	register_message = "Liebe/r Organisator/in von %s,\r\n\r\nvielen Dank für Deine Registrierung beim Mitmach-Republik e.V. Bitte bestätige Deine Registrierung, in dem Du auf den folgenden Link klickst:\r\n\r\n%s/approve/%s\r\n\r\nDas Team der Mitmach-Republik"
	password_subject = "Deine neue E-Mail-Adresse bei mitmachrepublik.de"
	password_message = "Liebe/r Organisator/in von %s,\r\n\r\nbitte bestätige Deine neue E-Mail-Adresse, in dem Du auf den folgenden Link klickst:\r\n\r\n%s/approve/%s\r\n\r\nDas Team der Mitmach-Republik"
	resetpwd_subject = "Neues Kennwort für mitmachrepublik.de"
	resetpwd_message = "Liebe/r Organisator/in von %s,\r\n\r\nDu hast auf unserer Webseite einen Link zur Neueingabe Deines Kennworts angefordet. Bitte klicke dazu auf den folgenden Link:\r\n\r\n%s/?auth=%s#newpwd\r\n\r\nDas Team der Mitmach-Republik"
	ga_dev           = "UA-61290824-1"
	ga_test          = "UA-61290824-2"
	ga_www           = "UA-61290824-3"
	google_api_key   = "AIzaSyAFzwmkGATzuHpcqV3g0yQEO77Vk66zXUM"
	facebook_app_id  = "138725613479008"
)

var (
	resultOK           = &appResult{Status: http.StatusOK}
	resultCreated      = &appResult{Status: http.StatusCreated}
	resultUnauthorized = &appResult{Status: http.StatusUnauthorized}
	resultBadRequest   = &appResult{Status: http.StatusBadRequest}
	resultNotFound     = &appResult{Status: http.StatusNotFound}
	resultConflict     = &appResult{Status: http.StatusConflict}

	sendAlertsService *SendAlertsService
)

func NewMmrApp(env string, host string, port int, tplDir, indexDir, imgServer, mongoUrl, dbName, smtpPass, fbAppSecret, fbUser, fbPass string) (*MmrApp, error) {

	database, err := NewMongoDb(mongoUrl, dbName)
	if err != nil {
		return nil, errors.New("init of MongoDB failed: " + err.Error())
	}

	users, err := NewUserService(database, "user")
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	events, err := NewEventService(database, "event", indexDir)
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	alerts, err := NewAlertService(database, "alert")
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	funcs := map[string]interface{}{
		"inc":                     inc,
		"dec":                     dec,
		"cut":                     cut,
		"intSlice":                intSlice,
		"dateFormat":              dateFormat,
		"timeFormat":              timeFormat,
		"datetimeFormat":          datetimeFormat,
		"longDatetimeFormat":      longDatetimeFormat,
		"iso8601Format":           iso8601Format,
		"noescape":                noescape,
		"strClip":                 strClip,
		"categoryIcon":            categoryIcon,
		"targetTitle":             targetTitle,
		"categoryTitle":           categoryTitle,
		"targetSearchUrl":         targetSearchUrl,
		"categorySearchUrl":       categorySearchUrl,
		"districtName":            districtName,
		"citypartName":            citypartName,
		"encodePath":              encodePath,
		"eventSearchUrl":          eventSearchUrl,
		"eventSearchUrlWithQuery": eventSearchUrlWithQuery,
		"organizerSearchUrl":      simpleOrganizerSearchUrl,
		"simpleEventSearchUrl":    simpleEventSearchUrl,
	}

	tpls, err := NewTemplates(tplDir+string(os.PathSeparator)+"*.tpl", funcs)
	if err != nil {
		return nil, errors.New("init of templates failed: " + err.Error())
	}

	emailAccount := &EmailAccount{"smtp.gmail.com", 465, "mitmachrepublik", smtpPass, &EmailAddress{"Mitmach-Republik e.V.", "mitmachrepublik@gmail.com"}}

	ga_code := ga_dev
	hostname := "http://dev.mitmachrepublik.de"
	if env == "www" {
		ga_code = ga_www
		hostname = "https://www.mitmachrepublik.de"
	} else if env == "test" {
		ga_code = ga_test
		hostname = "http://test.mitmachrepublik.de"
	}

	cities, err := events.Cities()
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	var adminId bson.ObjectId
	admin, err := users.LoadByEmail(ADMIN_EMAIL)
	if err == nil {
		adminId = admin.GetId()
	}

	services := make([]Service, 0)
	services = append(services, NewSessionService(3, emailAccount, database))
	services = append(services, NewScrapersService(3, emailAccount, events, adminId))
	services = append(services, NewUnusedImgService(4, emailAccount, database, imgServer))
	services = append(services, NewSendAlertsService(5, emailAccount, hostname, emailAccount, alerts))
	if env == "dev" {
		services = append(services, NewSpawnEventsService(12, emailAccount, database, events, imgServer))
	}

	return &MmrApp{host, port, tpls, imgServer, database, users, events, alerts, ga_code, hostname, emailAccount, NewLocationTree(cities), services, fbAppSecret, fbUser, fbPass}, nil
}

func (app *MmrApp) view(tpl string, w traffic.ResponseWriter, meta *metaTags, data bson.M) *appResult {

	if data == nil {
		data = bson.M{}
	}
	data["meta"] = meta
	data["ga_code"] = app.ga_code
	data["hostname"] = app.hostname
	data["districts"] = DistrictMap
	data["targets"] = TargetOrder
	data["targetMap"] = TargetMap
	data["categories"] = CategoryOrder
	data["categoryMap"] = CategoryMap
	data["googleApiKey"] = google_api_key

	var buf bytes.Buffer
	err := app.tpls.Execute(tpl, &buf, data)
	if err != nil {
		return &appResult{Status: http.StatusInternalServerError, Error: err}
	}

	w.Write(buf.Bytes())
	return resultOK
}

func (app *MmrApp) output(tpl string, w traffic.ResponseWriter, contentType string, meta *metaTags, data bson.M) *appResult {

	if data == nil {
		data = bson.M{}
	}
	data["meta"] = meta
	data["hostname"] = app.hostname

	var buf bytes.Buffer
	err := app.tpls.Execute(tpl, &buf, data)
	if err != nil {
		return &appResult{Status: http.StatusInternalServerError, Error: err}
	}

	w.Header().Set("Content-Type", contentType)
	w.Write(buf.Bytes())

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
		} else if !isEmpty(result.RedirectUrl) {
			w.Header().Set("Location", string(template.URL(result.RedirectUrl)))
			w.WriteHeader(http.StatusMovedPermanently)
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

func (app *MmrApp) sitemapPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		events, err := app.events.FindEvents()
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers, err := app.users.FindApproved()
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		i := 0
		topics := make([]string, len(Topics))
		for path := range Topics {
			topics[i] = path
			i++
		}

		return app.output("sitemap.tpl", w, "text/xml", nil, bson.M{"events": events, "organizers": organizers, "topics": topics})
	}()

	app.handle(w, result)
}

func (app *MmrApp) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 16
	query := ""
	place := ""
	dateIds := []int{TwoWeeks}
	timespans := timeSpans(dateIds)

	meta := metaTags{
		"Gemeinschaftliche Veranstaltungen zum Mitmachen!",
		"Gemeinschaftliche Veranstaltungen und Organisationen in Berlin und anderswo. Finde Veranstaltungen von Nachbarschaftszentren, Umweltverbänden, Bürgerinitiativen, Vereinen und Aktionsbündnissen für heute, morgen und am Wochenende in Deiner Umgebung.",
		"Gemeinsam aktiv werden.",
		app.hostname + "/images/mitmachrepublik.png",
		"Deine Seite für gemeinschaftliche Veranstaltungen und Organisationen. Finde Veranstaltungen von Nachbarschaftszentren, Umweltverbänden, Bürgerinitiativen und Vereinen in Deiner Umgebung. Mach mit bei gemeinsamen Projekten und Ideen!",
		true,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count("", place, timespans, nil, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		eventList := make([]*Event, 0)
		events := make(map[string][]*Event)
		for path, topic := range Topics {
			if topic.FrontPage == true {
				result, err := app.events.Search(query, topic.Place, timeSpans(topic.DateIds), topic.TargetIds, topic.CategoryIds, false, 0, pageSize)
				if err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}
				events[path] = result.Events
				eventList = append(eventList, result.Events...)
			}
		}

		organizers, err := app.users.FindForEvents(eventList)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		for topic, topicEvents := range events {
			approvedEvents := make([]*Event, 0)
			for _, event := range topicEvents {
				if organizers[event.OrganizerId].Approved {
					approvedEvents = append(approvedEvents, event)
				}
				if len(approvedEvents) == 8 {
					break
				}
			}
			events[topic] = approvedEvents
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": events2RSSItems(eventList)})
		} else {
			return app.view("start.tpl", w, &meta, bson.M{"topics": Topics, "events": events, "eventCnt": eventCnt, "organizerCnt": organizerCnt, "organizers": organizers, "timespans": timespans, "dates": DateOrder, "dateMap": DateIdMap})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) approvePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{
		"Registrierung bestätigen | Mitmach-Republik e.V.",
		"",
		"Registrierung bestätigen",
		app.hostname + "/images/mitmachrepublik.png",
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
			}
		} else {
			err = errors.New("No organizer id given.")
		}

		return app.view("approve.tpl", w, &meta, bson.M{"approved": err == nil})
	}()

	app.handle(w, result)
}

func (app *MmrApp) nlUnsubscribe(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{
		"Benachrichtigung abbestellen | Mitmach-Republik e.V.",
		"",
		"Benachrichtigung abbestellen",
		app.hostname + "/images/mitmachrepublik.png",
		"",
		false,
	}

	result := func() *appResult {

		var alert *Alert
		var err error = nil

		if bson.IsObjectIdHex(r.Param("id")) {
			alert, err = app.alerts.Load(bson.ObjectIdHex(r.Param("id")))
			if err == nil {
				err = app.alerts.Delete(alert.Id)
			}
		} else {
			err = errors.New("No email alert id given.")
		}

		return app.view("nl_unsubscribe.tpl", w, &meta, bson.M{"unsubscribed": err == nil})
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 10
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	if r.Param("fmt") == "RSS" {
		pageSize = 1000
		page = 0
	}

	place := r.Param("place")
	query := r.Param("query")
	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	var dateIds []int
	if r.Param("dateIds") == "aktuell" {
		dateIds = []int{FromNow}
	} else {
		dateIds = str2Int(strings.Split(r.Param("dateIds"), ","))
	}
	timespans := timeSpans(dateIds)
	targetIds := str2Int(strings.Split(r.Param("targetIds"), ","))
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	title := "Veranstaltungen"
	descr := title
	if len(targetIds) > 0 {
		if targetIds[0] == 0 {
			descr += " für " + strConcat(TargetOrder)
		} else {
			targets := make([]string, len(targetIds))
			for i, targetId := range targetIds {
				targets[i] = TargetIdMap[targetId]
			}
			title += " für " + strConcat(targets)
			descr += " für " + strConcat(targets)
		}
	}
	if len(categoryIds) > 0 {
		if categoryIds[0] == 0 {
			descr += " für " + strConcat(CategoryOrder)
		} else {
			categories := make([]string, len(categoryIds))
			for i, categoryId := range categoryIds {
				categories[i] = CategoryIdMap[categoryId]
			}
			if len(categories) == 1 {
				title += " aus der Kategorie " + strConcat(categories)
				descr += " aus der Kategorie " + strConcat(categories)
			} else {
				title += " aus den Kategorien " + strConcat(categories)
				descr += " aus den Kategorien " + strConcat(categories)
			}
		}
	}
	dateNames := ""
	if len(dateIds) > 0 && dateIds[0] != FromNow {
		dates := make([]string, len(dateIds))
		for i, dateId := range dateIds {
			dates[i] = DateIdMap[dateId]
		}
		dateNames = " " + strConcat(dates)
		title += dateNames
		descr += dateNames
	}
	if !isEmpty(place) {
		title += " in " + place
		descr += " in " + place
	}
	descr += "."

	meta := metaTags{
		title + " | Mitmach-Republik e.V.",
		descr,
		title,
		app.hostname + "/images/mitmachrepublik.png",
		descr,
		true,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count(query, place, timespans, targetIds, categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.events.Search(query, place, timespans, targetIds, categoryIds, false, page, pageSize)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers, err := app.users.FindForEvents(result.Events)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": events2RSSItems(result.Events)})
		} else {
			pageCount := pageCount(result.Count, pageSize)
			pages := make([]int, pageCount)
			for i := 0; i < pageCount; i++ {
				pages[i] = i
			}
			maxPage := pageCount - 1
			return app.view("events.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "organizers": organizers, "query": query, "place": place, "radius": radius, "timespans": timespans, "dates": DateOrder, "dateMap": DateIdMap, "dateIds": dateIds, "dateNames": dateNames, "targetIds": targetIds, "categoryIds": categoryIds, "noindex": true})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) landingPage(w traffic.ResponseWriter, r *traffic.Request) {

	path := r.Param("path")

	var topic *Topic
	for topicPath, topicData := range Topics {
		if path == topicPath {
			topic = &topicData
			break
		}
	}

	if topic == nil {
		app.handle(w, resultNotFound)
		return
	}

	query := ""

	pageSize := 10
	page, err := strconv.Atoi(r.Param("p"))
	if err != nil {
		page = 0
	}

	if r.Param("fmt") == "RSS" {
		pageSize = 1000
		page = 0
	}

	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	timespans := timeSpans(topic.DateIds)

	descr := "Veranstaltungen " + topic.Name
	if len(topic.CategoryIds) > 0 {
		if topic.CategoryIds[0] == 0 {
			descr += " für " + strConcat(CategoryOrder)
		} else {
			categories := make([]string, len(topic.CategoryIds))
			for i, categoryId := range topic.CategoryIds {
				categories[i] = CategoryIdMap[categoryId]
			}
			if len(categories) == 1 {
				descr += " aus der Kategorie " + strConcat(categories)
			} else {
				descr += " aus den Kategorien " + strConcat(categories)
			}
		}
	}

	title := "Veranstaltungen für " + topic.Name
	if len(topic.Place) > 0 {
		title += " in " + topic.Place
		descr += " in " + topic.Place
	}

	meta := metaTags{
		title + " | Mitmach-Republik e.V.",
		descr,
		title,
		app.hostname + "/images/mitmachrepublik.png",
		descr,
		true,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count(query, topic.Place, timespans, topic.TargetIds, topic.CategoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(topic.Place, topic.CategoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.events.Search(query, topic.Place, timespans, topic.TargetIds, topic.CategoryIds, false, page, pageSize)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers, err := app.users.FindForEvents(result.Events)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": events2RSSItems(result.Events)})
		} else {
			pageCount := pageCount(result.Count, pageSize)
			pages := make([]int, pageCount)
			for i := 0; i < pageCount; i++ {
				pages[i] = i
			}
			maxPage := pageCount - 1
			return app.view("events.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "organizers": organizers, "query": query, "place": topic.Place, "radius": radius, "timespans": timespans, "dates": DateOrder, "dateMap": DateIdMap, "dateIds": topic.DateIds, "targetIds": topic.TargetIds, "categoryIds": topic.CategoryIds, "headline": title, "noindex": len(r.Param("p")) > 0, "altPage": true})
		}
	}()

	app.handle(w, result)
}

func (app *MmrApp) nlEventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	alertId := r.Param("id")

	query := r.Param("query")
	place := r.Param("place")
	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	dateIds := str2Int(strings.Split(r.Param("dateIds"), ","))
	timespans := timeSpans(dateIds)

	targetIds := str2Int(strings.Split(r.Param("targetIds"), ","))
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	result := func() *appResult {

		result, err := app.events.Search(query, place, timespans, targetIds, categoryIds, false, 0, 1000)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		if result.Count == 0 {
			return resultNotFound
		}

		organizers, err := app.users.FindForEvents(result.Events)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return app.output("nl_events.tpl", w, "text/html", nil, bson.M{"alertId": alertId, "results": result.Count, "events": result.Events, "organizers": organizers, "place": place, "radius": radius, "dateIds": dateIds, "timespans": timespans, "targetIds": targetIds, "categoryIds": categoryIds})
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventPage(w traffic.ResponseWriter, r *traffic.Request) {

	radius := 2
	dateIds := []int{TwoWeeks}
	from, _ := strconv.ParseInt(r.Param("from"), 10, 64)
	if from == 0 {
		from = time.Now().Unix()
	}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		event, err := app.events.Load(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		duration := time.Duration(0)
		if !event.End.IsZero() {
			duration = event.End.Sub(event.Start)
		}
		start := event.NextDate(time.Unix(from, 0))
		end := start.Add(duration)

		organizer, err := app.users.Load(event.OrganizerId)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		place := citypartName(event.Addr)

		eventCnt, err := app.events.Count("", place, timeSpans(dateIds), nil, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		imageUrl := ""
		if !isEmpty(event.Image) {
			imageUrl = app.hostname + "/bild/" + event.Image
		}

		title := event.Title
		if !isEmpty(place) {
			title += " in " + place
		}
		if !isEmpty(event.Addr.Name) {
			title += " (" + event.Addr.Name + ")"
		}

		meta := metaTags{
			title + " am " + dateFormat(event.Start) + " | Mitmach-Republik e.V.",
			strClip(event.PlainDescription(), 160),
			title,
			imageUrl,
			strClip(event.PlainDescription(), 160),
			false,
		}

		similiars, err := app.events.FindSimilarEvents(event, 8)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers, err := app.users.FindForEvents(similiars)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return app.view("event.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "place": place, "radius": radius, "event": event, "from": time.Unix(from, 0), "start": start, "end": end, "organizer": organizer, "similiars": similiars, "organizers": organizers})
	}()

	app.handle(w, result)
}

func (app *MmrApp) sendEventPage(w traffic.ResponseWriter, r *traffic.Request) {

	from, _ := strconv.ParseInt(r.Param("from"), 10, 64)
	if from == 0 {
		from = time.Now().Unix()
	}

	meta := metaTags{"Veranstaltung empfehlen | Mitmach-Republik e.V.", "", "", "", "", false}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		event, err := app.events.Load(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		return app.view("sendevent.tpl", w, &meta, bson.M{"event": &event, "from": from})
	}()

	app.handle(w, result)
}

func (app *MmrApp) emailAlertPage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"E-Mail-Benachrichtung | Mitmach-Republik e.V.", "", "", "", "", false}

	result := func() *appResult {

		data := bson.M{}
		for key, value := range r.Params() {
			if len(value) > 0 {
				data[key] = value[0]
			}
		}

		return app.view("emailalert.tpl", w, &meta, data)
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizersPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 10
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

	title := "Organisatoren"
	descr := title
	if len(categoryIds) > 0 {
		if categoryIds[0] == 0 {
			descr += " für " + strConcat(CategoryOrder)
		} else {
			categories := make([]string, len(categoryIds))
			for i, categoryId := range categoryIds {
				categories[i] = CategoryIdMap[categoryId]
			}
			title += " für " + strConcat(categories)
			descr += " für " + strConcat(categories)
		}
	}
	if !isEmpty(place) {
		title += " in " + place
		descr += " in " + place
	}
	descr += "."

	meta := metaTags{
		title + " | Mitmach-Republik e.V.",
		descr,
		title,
		app.hostname + "/images/mitmachrepublik.png",
		descr,
		false,
	}

	result := func() *appResult {

		eventCnt, err := app.events.Count("", place, timeSpans(dateIds), nil, categoryIds)
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

		return app.view("organizers.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "organizers": result.Organizers, "place": place, "radius": radius, "categoryIds": categoryIds, "noindex": true})
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizerPage(w traffic.ResponseWriter, r *traffic.Request) {

	pageSize := 10
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
	timespans := timeSpans(dateIds)

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		organizer, err := app.users.Load(bson.ObjectIdHex(r.Param("id")))
		if err != nil {
			return resultNotFound
		}

		place := citypartName(organizer.Addr)

		eventCnt, err := app.events.Count("", place, timespans, nil, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.users.Count(place, nil)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		result, err := app.events.SearchFutureEventsOfUser(organizer.Id, page, pageSize)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers := map[bson.ObjectId]*User{organizer.Id: organizer}

		imageUrl := ""
		if !isEmpty(organizer.Image) {
			imageUrl = app.hostname + "/bild/" + organizer.Image
		}
		title := "Veranstaltungen"
		if organizer.Name != "Mitmach-Republik e.V." {
			title += " von " + organizer.Name
			if !isEmpty(place) {
				title += " aus " + place
			}
		}
		meta := metaTags{
			title + " | Mitmach-Republik e.V.",
			strClip(organizer.PlainDescription(), 160),
			title,
			imageUrl,
			strClip(organizer.PlainDescription(), 160),
			true,
		}

		if r.Param("fmt") == "RSS" {
			return app.output("rss.tpl", w, "application/rss+xml", &meta, bson.M{"items": events2RSSItems(result.Events)})
		} else {
			pageCount := pageCount(result.Count, pageSize)
			pages := make([]int, pageCount)
			for i := 0; i < pageCount; i++ {
				pages[i] = i
			}
			maxPage := pageCount - 1
			return app.view("organizer.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "organizers": organizers, "place": place, "radius": radius, "timespans": timespans, "organizer": organizer, "noindex": page > 0})
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

	pageSize := 10
	page, err := strconv.Atoi(r.Param("page"))
	if err != nil {
		page = 0
	}

	search := strings.Trim(r.Param("query"), " ")
	meta := metaTags{"Verwaltung | Mitmach-Republik e.V.", "", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		result, err := app.events.SearchEventsOfUser(user.Id, search, page, pageSize, "-start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizers := map[bson.ObjectId]*User{user.Id: user}

		pageCount := result.Count / pageSize
		if pageCount == 0 {
			pageCount = 1
		}
		pages := make([]int, pageCount)
		for i := 0; i < pageCount; i++ {
			pages[i] = i
		}
		maxPage := pageCount - 1

		return app.view("admin.tpl", w, &meta, bson.M{"user": user, "query": search, "results": result.Count, "organizers": organizers, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "timespans": [][]time.Time{[]time.Time{time.Time{}, time.Time{}}}})
	}()

	app.handle(w, result)
}

func (app *MmrApp) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"Profil bearbeiten | Mitmach-Republik e.V.", "", "", "", "", false}

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

	meta := metaTags{"E-Mail-Adresse oder Kennwort ändern | Mitmach-Republik e.V.", "", "", "", "", false}

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

	meta := metaTags{"Veranstaltung bearbeiten | Mitmach-Republik e.V.", "", "", "", "", false}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		organizers := make(map[bson.ObjectId]string)
		if user.IsAdmin() {
			users, err := app.users.FindUsers()
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}
			for _, user := range users {
				organizers[user.Id] = user.Name
			}
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
			event.ImageCredit = oldEvent.ImageCredit
			event.Targets = oldEvent.Targets
			event.Categories = oldEvent.Categories
			event.Descr = oldEvent.Descr
			event.Addr = oldEvent.Addr
			event.Web = oldEvent.Web
			event.Recurrency = oldEvent.Recurrency
			event.RecurrencyEnd = oldEvent.RecurrencyEnd
			event.Monthly = oldEvent.Monthly
			event.Weekly = oldEvent.Weekly
		} else {
			event = new(Event)
		}

		return app.view("event_edit.tpl", w, &meta, bson.M{"user": user, "event": event, "organizers": organizers})
	}()

	app.handle(w, result)
}

func (app *MmrApp) staticPage(w traffic.ResponseWriter, template, headline string) {

	meta := metaTags{headline + " | Mitmach-Republik e.V.", headline, "", "", "", false}

	result := func() *appResult {
		return app.view(template, w, &meta, nil)
	}()

	app.handle(w, result)
}

func (app *MmrApp) searchHandler(w traffic.ResponseWriter, r *traffic.Request) {

	query := strings.Trim(r.PostFormValue("query"), " ")
	if isEmpty(query) {
		query = strings.Trim(r.PostFormValue("fulltextsearch"), " ")
	}
	place := app.locations.Normalize(strings.Trim(r.PostFormValue("place"), " "))

	radius, err := strconv.Atoi(r.PostFormValue("radius"))
	if err != nil {
		radius = 0
	}

	targetIds := str2Int(r.Form["target"])
	if len(targetIds) == 0 {
		targetIds = append(targetIds, 0)
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
		if isEmpty(query) {
			path = "/veranstaltungen/" + eventSearchUrl(place, targetIds, categoryIds, dateIds, radius)
		} else {
			path = "/veranstaltungen/" + eventSearchUrlWithQuery(place, targetIds, categoryIds, dateIds, radius, query)
		}
	}

	w.Header().Set("Location", path+"#events")
	w.WriteHeader(http.StatusFound)
}

func (app *MmrApp) uploadHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		file, _, err := r.FormFile("file")
		if err != nil {
			return resultBadRequest
		}

		filename := uuid.New().String() + ".jpg"
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

	return SendEmail(app.emailAccount, to, replyTo, subject, "text/plain", body)
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
		user.ImageCredit = data.ImageCredit
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
			} else {
				event.Source = oldEvent.Source
				event.SourceId = oldEvent.SourceId
				event.FacebookId = oldEvent.FacebookId
			}
		} else {
			created = true
			event.Id = bson.NewObjectId()
		}

		if !user.IsAdmin() || !event.OrganizerId.Valid() {
			event.OrganizerId = user.Id
		}

		organizer := user
		if organizer.Id != event.OrganizerId {
			organizer, err = app.users.Load(event.OrganizerId)
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}
		}

		err = app.events.Store(event)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		go func() {
			app.locations.Add(event.Addr.City)
			if app.fbAppSecret != "" && app.fbUser != "" && app.fbPass != "" {
				if event.Facebook == true && event.FacebookId == "" {
					client, err := NewFacebookClient(app.hostname, facebook_app_id, app.fbAppSecret, app.fbUser, app.fbPass)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}

					event.FacebookId, err = client.PostEvent(event)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}

					err = app.events.Store(event)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}

				} else if event.Facebook == false && event.FacebookId != "" {
					client, err := NewFacebookClient(app.hostname, facebook_app_id, app.fbAppSecret, app.fbUser, app.fbPass)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}

					err = client.DeletePost(event.FacebookId)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}

					event.FacebookId = ""
					err = app.events.Store(event)
					if err != nil {
						traffic.Logger().Print(err.Error())
						return
					}
				}
			}
		}()

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

func (app *MmrApp) sendPasswordMailHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		form, err := (&Request{r}).ReadSendPassword()
		if err != nil {
			return resultBadRequest
		}

		user, err := app.users.LoadByEmail(form.Email)
		if err != nil {
			// we don't tell if the email was not found
			return resultOK
		}

		sessionId, err := app.database.CreateSession(user.GetId())
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = app.sendEmail(app.emailAccount.From, &EmailAddress{user.Name, user.Email}, resetpwd_subject, fmt.Sprintf(resetpwd_message, user.Name, app.hostname, sessionId.Hex()))
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) emailAlertHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		var alert Alert
		err := (&Request{r}).ReadJson(&alert)
		if err != nil {
			return resultBadRequest
		}

		alert.Id = bson.NewObjectId()
		err = app.alerts.Store(&alert)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	app.handle(w, result)
}

func (app *MmrApp) typeAheadHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		events, err := app.events.SearchText(strings.Trim(r.Param("query"), " "))
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		} else {
			return &appResult{Status: http.StatusOK, JSON: events}
		}
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

	dateIds := str2Int(strings.Split(r.Param("dateIds"), ","))
	targetIds := str2Int(strings.Split(r.Param("targetIds"), ","))
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))
	query := strings.Trim(r.Param("query"), " ")
	place := app.locations.Normalize(strings.Trim(r.Param("place"), " "))

	result := func() *appResult {

		cnt, err := app.events.Count(query, place, timeSpans(dateIds), targetIds, categoryIds)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusOK, JSON: cnt}
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizerCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))
	place := app.locations.Normalize(strings.Trim(r.Param("place"), " "))

	result := func() *appResult {

		cnt, err := app.users.Count(place, categoryIds)
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

	router.Get("/veranstaltungen/:place/:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen//:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen/:place/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen//:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen/:place/:dateIds/:categoryIds/:radius/:categories/:page", app.eventsPage)
	router.Get("/veranstaltungen//:dateIds/:categoryIds/:radius/:categories/:page", app.eventsPage)

	router.Get("/newsletter/veranstaltungen/:place/:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:id", app.nlEventsPage)
	router.Get("/newsletter/veranstaltungen//:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:id", app.nlEventsPage)

	router.Get("/veranstaltung/:place/:targets/:categories/:dateId/:id/:title", app.eventPage)
	router.Get("/veranstaltung//:targets/:categories/:dateId/:id/:title", app.eventPage)
	router.Get("/veranstaltung/:place//:categories/:dateId/:id/:title", app.eventPage)
	router.Get("/veranstaltung///:categories/:dateId/:id/:title", app.eventPage)

	router.Get("/veranstaltung/:place/:targets/:categories/:id/:title", app.eventPage)
	router.Get("/veranstaltung//:targets/:categories/:id/:title", app.eventPage)
	router.Get("/veranstaltung/:place//:categories/:id/:title", app.eventPage)
	router.Get("/veranstaltung///:categories/:id/:title", app.eventPage)

	router.Get("/veranstalter/:place/:categoryIds/:categories/:page", app.organizersPage)
	router.Get("/veranstalter//:categoryIds/:categories/:page", app.organizersPage)
	router.Get("/veranstalter/:id/:title/:page", app.organizerPage)

	router.Get("/sitemap.xml", app.sitemapPage)
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
	router.Get("/wir-ueber-uns", func(w traffic.ResponseWriter, r *traffic.Request) {
		app.staticPage(w, "aboutus.tpl", "Wir über uns")
	})

	router.Get("/dialog/contact", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "contact.tpl", "") })
	router.Get("/dialog/registered", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "registered.tpl", "") })
	router.Get("/dialog/login", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "login.tpl", "") })
	router.Get("/dialog/register", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "register.tpl", "") })
	router.Get("/dialog/password", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, "password_reset.tpl", "") })
	router.Get("/dialog/sendevent/:id", app.sendEventPage)
	router.Get("/dialog/emailalert/:place/:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.emailAlertPage)
	router.Get("/dialog/emailalert//:dates/:dateIds/:targetIds/:categoryIds/:radius/:targets/:categories/:page", app.emailAlertPage)

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
	router.Post("/sendpasswordmail", app.sendPasswordMailHandler)
	router.Post("/emailalert", app.emailAlertHandler)
	router.Get("/newsletter/unsubscribe/:id", app.nlUnsubscribe)

	router.Get("/typeahead/:query", app.typeAheadHandler)
	router.Get("/location/:location", app.locationHandler)
	router.Get("/eventcount/:query/:place/:dateIds/:targetIds/:categoryIds", app.eventCountHandler)
	router.Get("/eventcount/:query//:dateIds/:targetIds/:categoryIds", app.eventCountHandler)
	router.Get("/eventcount//:place/:dateIds/:targetIds/:categoryIds", app.eventCountHandler)
	router.Get("/eventcount///:dateIds/:targetIds/:categoryIds", app.eventCountHandler)
	router.Get("/organizercount/:place/:categoryIds", app.organizerCountHandler)
	router.Get("/organizercount//:categoryIds", app.organizerCountHandler)

	router.Post("/login", app.loginHandler)
	router.Post("/logout", app.logoutHandler)

	router.Get("/approve/:id", app.approvePage)
	router.Delete("/event/:id", app.deleteEventHandler)

	router.Get("/:path", app.landingPage)

	router.Run()
}

func (app *MmrApp) RunScrapers() error {

	var adminId bson.ObjectId
	admin, err := app.users.LoadByEmail(ADMIN_EMAIL)
	if err == nil {
		adminId = admin.Id
	}
	scrapers := NewScrapersService(0, app.emailAccount, app.events, adminId)
	return scrapers.Run()
}

func (app *MmrApp) Stop() error {
	return app.events.Stop()
}
