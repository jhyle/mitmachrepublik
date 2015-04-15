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
	"time"
)

type (
	MmrApp struct {
		host         string
		port         int
		tpls         *Templates
		imgServer    string
		database     Database
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
		JSON   interface{}
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

	err = database.Table("user").EnsureIndices("name", "email", "approved", "categories", "addr.city", "addr.pcode")
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	err = database.Table("event").EnsureIndices("organizerid", "approved", "start", "categories", "addr.city", "addr.pcode")
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
		"strClip":            strClip,
		"categoryIcon":       categoryIcon,
		"categoryTitle":      categoryTitle,
		"categorySearchUrl":  categorySearchUrl,
		"districtName":       districtName,
		"citypartName":       citypartName,
		"encodePath":         encodePath,
		"eventUrl":           eventUrl,
		"eventSearchUrl":     simpleEventSearchUrl,
		"organizerUrl":       organizerUrl,
		"organizerSearchUrl": simpleOrganizerSearchUrl,
	}

	tpls, err := NewTemplates(tplDir+string(os.PathSeparator)+"*.tpl", funcs)
	if err != nil {
		return nil, errors.New("init of templates failed: " + err.Error())
	}

	emailAccount := &EmailAccount{"smtp.gmail.com", 587, "mitmachrepublik", "mitmachen", &EmailAddress{"Mitmach-Republik", "mitmachrepublik@gmail.com"}}

	ga_code := ga_dev
	hostname := "dev.mitmachrepublik.de"
	if env == "www" {
		ga_code = ga_www
		hostname = "www.mitmachrepublik.de"
	} else if env == "test" {
		ga_code = ga_test
		hostname = "test.mitmachrepublik.de"
	}

	var cities []string
	err = database.Table("event").Distinct(bson.M{}, "addr.city", &cities)
	if err != nil {
		return nil, errors.New("init of database failed: " + err.Error())
	}

	services := make([]Service, 0, 3)
	services = append(services, NewSessionService(60, database))
	services = append(services, NewUnusedImgService(3600, database, imgServer))
	if env != "www" {
		services = append(services, NewSpawnEventsService(3600, database, imgServer))
	}

	return &MmrApp{host, port, tpls, imgServer, database, ga_code, hostname, emailAccount, NewLocationTree(cities), services}, nil
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
		}
	}

	if result.JSON != nil && !w.BodyWritten() {
		w.WriteJSON(result.JSON)
	}
}

func (app *MmrApp) countEvents(place string, categoryIds []int, dateNames []string) (int, error) {

	query := buildQuery(place, timeSpans(dateNames), categoryIds)
	return app.database.Table("event").Count(query)
}

func (app *MmrApp) countOrganizers(place string) (int, error) {

	query := buildQuery(place, nil, nil)
	return app.database.Table("user").Count(query)
}

func (app *MmrApp) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	eventsPerRow := 4
	place := "Berlin"
	dateNames := []string{"aktuell"}

	meta := metaTags{
		"Willkommen in der Mitmach-Republik!",
		"Gemeinsam aktiv werden.",
		"http://" + app.hostname + "/images/mitmachrepublik.gif",
		"Hier findest Du Veranstaltungen und Organisationen zum Mitmachen. Suche nach Nachbarschaftstreffen, Sportvereinen, gemeinnützigen Initiativen, religiösen Gemeinden und anderen Vereinen in Deiner Umgebung. Mach mit bei gemeinsamen Projekten und Ideen.",
	}

	result := func() *appResult {

		eventCnt, err := app.countEvents(place, nil, dateNames)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.countOrganizers(place)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		var result EventSearchResult
		err = app.database.Table("event").Search(buildQuery(place, timeSpans(dateNames), nil), 0, eventsPerRow*2, &result, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		var events [][]*Event
		if len(result.Events) > 0 {
			events = make([][]*Event, ((len(result.Events)-1)/eventsPerRow)+1)
			for i := range events {

				rowSize := len(result.Events) - i*eventsPerRow
				if rowSize > eventsPerRow {
					rowSize = eventsPerRow
				}
				events[i] = make([]*Event, rowSize)

				for j := 0; j < rowSize; j++ {
					events[i][j] = result.Events[i*eventsPerRow+j]
				}
			}
		} else {
			events = make([][]*Event, 0)
		}

		return app.view("start.tpl", w, &meta, bson.M{"events": events, "eventCnt": eventCnt, "organizerCnt": organizerCnt})
	}()

	app.handle(w, result)
}

func (app *MmrApp) approvePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{
		"Registrierung bestätigen - Mitmach-Republik",
		"Registrierung bestätigen",
		"http://" + app.hostname + "/images/mitmachrepublik.gif",
		"",
	}

	result := func() *appResult {

		var err error = nil

		if bson.IsObjectIdHex(r.Param("id")) {
			var user User
			userId := bson.ObjectIdHex(r.Param("id"))

			if err = app.database.Table("user").LoadById(userId, &user); err == nil {

				user.Approved = true
				if _, err = app.database.Table("user").UpsertById(user.Id, &user); err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}

				var events []Event
				if err = app.database.Table("event").Find(bson.M{"organizerid": user.Id}, &events, "_id"); err == nil {
					for _, event := range events {
						event.Approved = true
						_, err = app.database.Table("event").UpsertById(event.Id, &event)
						if err != nil {
							break
						}
					}
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

	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	place := r.Param("place")
	dateNames := strings.Split(r.Param("dates"), ",")
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	meta := metaTags{
		"Veranstaltungen in " + place + " - Mitmach-Republik",
		"Veranstaltungen in " + place,
		"http://" + app.hostname + "/images/mitmachrepublik.gif",
		"",
	}

	result := func() *appResult {

		eventCnt, err := app.countEvents(place, categoryIds, dateNames)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.countOrganizers(place)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		var result EventSearchResult
		query := buildQuery(place, timeSpans(dateNames), categoryIds)
		err = app.database.Table("event").Search(query, page*pageSize, pageSize, &result, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerNames := make(map[bson.ObjectId]string)
		for _, event := range result.Events {
			if _, found := organizerNames[event.OrganizerId]; !found {
				var user User
				err = app.database.Table("user").LoadById(event.OrganizerId, &user)
				if err != nil {
					return &appResult{Status: http.StatusInternalServerError, Error: err}
				}
				organizerNames[event.OrganizerId] = user.Name
			}
		}

		pageCount := pageCount(result.Count, pageSize)
		pages := make([]int, pageCount)
		for i := 0; i < pageCount; i++ {
			pages[i] = i
		}
		maxPage := pageCount - 1

		return app.view("events.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "organizerNames": organizerNames, "place": place, "radius": radius, "dates": dateNames, "categoryIds": categoryIds})
	}()

	app.handle(w, result)
}

func (app *MmrApp) eventPage(w traffic.ResponseWriter, r *traffic.Request) {

	radius := 2
	dateNames := []string{"aktuell"}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		var event Event
		err := app.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), &event)
		if err != nil {
			return resultNotFound
		}

		var organizer User
		err = app.database.Table("user").LoadById(event.OrganizerId, &organizer)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		place := citypartName(event.Addr)

		eventCnt, err := app.countEvents(place, nil, dateNames)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.countOrganizers(place)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		imageUrl := ""
		if !isEmpty(event.Image) {
			imageUrl = "http://" + app.hostname + "/bild/" + event.Image
		}
		meta := metaTags{
			event.Title + " in " + place + " - Mitmach-Republik",
			event.Title + " in " + place,
			imageUrl,
			event.Descr,
		}

		return app.view("event.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "place": place, "radius": radius, "event": &event, "organizer": &organizer})
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
	dateNames := []string{"aktuell"}
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	meta := metaTags{
		"Veranstalter in " + place + " - Mitmach-Republik",
		"Veranstalter in " + place,
		"http://" + app.hostname + "/images/mitmachrepublik.gif",
		"",
	}

	result := func() *appResult {

		eventCnt, err := app.countEvents(place, categoryIds, dateNames)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.countOrganizers(place)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		var result OrganizerSearchResult
		query := buildQuery(place, nil, categoryIds)
		err = app.database.Table("user").Search(query, page*pageSize, pageSize, &result, "name")
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

	radius := 2
	dateNames := []string{"aktuell"}

	result := func() *appResult {

		if !bson.IsObjectIdHex(r.Param("id")) {
			return resultNotFound
		}

		var organizer User
		err := app.database.Table("user").LoadById(bson.ObjectIdHex(r.Param("id")), &organizer)
		if err != nil {
			return resultNotFound
		}

		place := citypartName(organizer.Addr)

		eventCnt, err := app.countEvents(place, nil, dateNames)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := app.countOrganizers(place)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		var result EventSearchResult
		query := bson.M{"$and": []bson.M{bson.M{"organizerid": organizer.Id}, bson.M{"approved": true}, bson.M{"start": bson.M{"$gte": time.Now()}}}}
		err = app.database.Table("event").Search(query, page*pageSize, pageSize, &result, "start")
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		pageCount := pageCount(result.Count, pageSize)
		pages := make([]int, pageCount)
		for i := 0; i < pageCount; i++ {
			pages[i] = i
		}
		maxPage := pageCount - 1

		organizerNames := map[bson.ObjectId]string{organizer.Id: organizer.Name}

		imageUrl := ""
		if !isEmpty(organizer.Image) {
			imageUrl = "http://" + app.hostname + "/bild/" + organizer.Image
		}
		name := organizer.Name
		if !isEmpty(place) {
			name += " aus " + place
		}
		meta := metaTags{
			name + " - Mitmach-Republik",
			name,
			imageUrl,
			organizer.Descr,
		}

		return app.view("organizer.tpl", w, &meta, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "results": result.Count, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events, "organizerNames": organizerNames, "place": place, "radius": radius, "organizer": organizer})
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

	meta := metaTags{"Verwaltung - Mitmach-Republik", "", "", ""}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		var result EventSearchResult
		err = app.database.Table("event").Search(bson.M{"organizerid": user.Id}, page*pageSize, pageSize, &result, "-start")
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

		return app.view("admin.tpl", w, &meta, bson.M{"user": user, "page": page, "pages": pages, "maxPage": maxPage, "events": result.Events})
	}()

	app.handle(w, result)
}

func (app *MmrApp) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	meta := metaTags{"Profil bearbeiten - Mitmach-Republik", "", "", ""}

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

	meta := metaTags{"E-Mail-Adresse oder Kennwort ändern - Mitmach-Republik", "", "", ""}

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

	meta := metaTags{"Veranstaltung bearbeiten - Mitmach-Republik", "", "", ""}

	result := func() *appResult {

		user, err := app.checkSession((&Request{r}))
		if err != nil {
			return resultUnauthorized
		}

		var data Event
		event := &data
		if bson.IsObjectIdHex(r.Param("id")) {

			err = app.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
			if err != nil {
				return &appResult{Status: http.StatusNotFound, Error: err}
			}

			if event.OrganizerId != user.Id {
				return resultBadRequest
			}
		}

		return app.view("event_edit.tpl", w, &meta, bson.M{"user": user, "event": event})
	}()

	app.handle(w, result)
}

func (app *MmrApp) staticPage(w traffic.ResponseWriter, r *traffic.Request, template, headline string) {

	meta := metaTags{headline + " - Mitmach-Republik", headline, "", ""}

	result := func() *appResult {
		return app.view(template, w, &meta, nil)
	}()

	app.handle(w, result)
}

func (app *MmrApp) searchHandler(w traffic.ResponseWriter, r *traffic.Request) {

	place := strings.Trim(r.PostFormValue("place"), " ")
	if isEmpty(place) {
		place = "Berlin"
	}

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

		err = validateUser(app.database, user)
		if err != nil {
			return resultConflict
		}

		user.Id = bson.NewObjectId()
		_, err = app.database.Table("user").UpsertById(user.Id, user)
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

		_, err = app.database.Table("user").UpsertById(user.Id, user)
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
			err := app.sendEmail(&EmailAddress{user.Name, user.Email}, nil, password_subject, fmt.Sprintf(password_message, app.hostname, user.Name, user.Id.Hex()))
			if err != nil {
				return &appResult{Status: http.StatusInternalServerError, Error: err}
			}
		}

		_, err = app.database.Table("user").UpsertById(user.Id, user)
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

		err = app.database.Table("event").Delete(bson.M{"organizerid": user.Id})
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = app.database.Table("user").DeleteById(user.Id)
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
			var oldData Event
			oldEvent := &oldData
			app.database.Table("event").LoadById(event.Id, oldEvent)
			if oldEvent.OrganizerId != user.Id {
				return resultBadRequest
			}
		} else {
			created = true
			event.Id = bson.NewObjectId()
		}
		event.OrganizerId = user.Id
		event.Approved = user.Approved

		_, err = app.database.Table("event").UpsertById(event.Id, event)
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

		categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

		dateIds := str2Int(strings.Split(r.Param("dateIds"), ","))
		dates := make([]string, len(dateIds))
		for i, dateId := range dateIds {
			dates[i] = DateIdMap[dateId]
		}

		cnt, err := app.countEvents(r.Param("place"), categoryIds, dates)
		if err != nil {
			return &appResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &appResult{Status: http.StatusOK, JSON: cnt}
	}()

	app.handle(w, result)
}

func (app *MmrApp) organizerCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *appResult {

		cnt, err := app.countOrganizers(r.Param("place"))
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

		var data Event
		event := &data
		err = app.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
		if err != nil {
			return resultOK
		}

		if event.OrganizerId != user.Id {
			return resultBadRequest
		}

		err = app.database.Table("event").DeleteById(event.Id)
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
	router.Get("/veranstaltungen/:place/:dates/:categoryIds/:radius/:categories/:page", app.eventsPage)
	router.Get("/veranstaltung//:date/:id/:title", app.eventPage)
	router.Get("/veranstaltung/:categories/:date/:id/:title", app.eventPage)
	router.Get("/veranstalter/:id/:title/:page", app.organizerPage)
	router.Get("/veranstalter/:place/:categoryIds/:categories/:page", app.organizersPage)

	router.Get("/impressum", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, r, "impressum.tpl", "Impressum") })
	router.Get("/disclaimer", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, r, "disclaimer.tpl", "Haftungsausschluss (Disclaimer)") })
	router.Get("/datenschutz", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, r, "datenschutz.tpl", "Datenschutzerklärung") })
	router.Get("/agbs", func(w traffic.ResponseWriter, r *traffic.Request) { app.staticPage(w, r, "agbs.tpl", "Allgemeine Geschäftsbedingungen") })

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
	router.Get("/organizercount/:place/:categoryIds", app.organizerCountHandler)
	router.Post("/login", app.loginHandler)
	router.Post("/logout", app.logoutHandler)

	router.Get("/approve/:id", app.approvePage)
	router.Delete("/event/:id", app.deleteEventHandler)

	router.Run()
}
