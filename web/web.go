package mmr

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/pilu/traffic"
	"labix.org/v2/mgo/bson"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type (
	WebServer struct {
		host         string
		port         int
		niceExpr     *regexp.Regexp
		tpls         *Templates
		imgServer    string
		database     Database
		emailAccount *EmailAccount
		locations    *LocationTree
	}

	emailAndPwd struct {
		Email string
		Pwd   string
	}

	webResult struct {
		Status int
		Error  error
		JSON   interface{}
	}
)

const (
	register_subject = "Deine Registrierung bei mitmachrepublik.de"
	register_message = "Liebe/r Organisator/in von %s,\r\n\r\nvielen Dank für die Registrierung bei der MitmachRepublik. Bitte bestätige Deine Registrierung, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://dev.mitmachrepublik.de/approve/%s\r\n\r\nDas Team der MitmachRepublik"
	password_subject = "Deine neue E-Mail-Adresse bei mitmachrepublik.de"
	password_message = "Liebe/r Organisator/in von %s,\r\n\r\nbitte bestätige Deine neue E-Mail-Adresse, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://dev.mitmachrepublik.de/approve/%s\r\n\r\nDas Team der MitmachRepublik"
)

var (
	resultOK         = &webResult{Status: http.StatusOK}
	resultCreated    = &webResult{Status: http.StatusCreated}
	resultBadRequest = &webResult{Status: http.StatusBadRequest}
	resultNotFound   = &webResult{Status: http.StatusNotFound}
	resultConflict   = &webResult{Status: http.StatusConflict}
)

func NewWebServer(host string, port int, tplDir, imgServer, mongoUrl, dbName string) (*WebServer, error) {

	niceExpr, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		return nil, err
	}

	database, err := NewMongoDb(mongoUrl, dbName)
	if err != nil {
		return nil, err
	}

	err = database.Table("user").EnsureIndices("email")
	if err != nil {
		return nil, err
	}

	err = database.Table("event").EnsureIndices("organizerid", "start")
	if err != nil {
		return nil, err
	}

	tpls, err := NewTemplates(tplDir + string(os.PathSeparator) + "*.tpl")
	if err != nil {
		return nil, err
	}

	emailAccount := &EmailAccount{"smtp.gmail.com", 587, "mitmachrepublik", "mitmachen", "MitmachRepublik <mitmachrepublik@gmail.com>"}

	var cities []string
	err = database.Table("event").Distinct(bson.M{}, "addr.city", &cities)
	if err != nil {
		return nil, err
	}

	return &WebServer{host, port, niceExpr, tpls, imgServer, database, emailAccount, NewLocationTree(cities)}, nil
}

func (web *WebServer) view(tpl string, w traffic.ResponseWriter, data bson.M) *webResult {

	err := web.tpls.Execute(tpl, w, data)
	if err != nil {
		return &webResult{Status: http.StatusInternalServerError, Error: err}
	}

	return resultOK
}

func (web *WebServer) handle(w traffic.ResponseWriter, result *webResult) {

	if result.Error != nil {
		traffic.Logger().Print(result.Error.Error())
	}

	if !w.Written() {
		w.WriteHeader(result.Status)
	}

	if result.JSON != nil && !w.BodyWritten() {
		w.WriteJSON(result.JSON)
	}
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
		placesQuery := make([]bson.M, len(postcodes) + 1)
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

func (web *WebServer) countEvents(place string, categoryIds []int, dateNames []string) (int, error) {

	query := buildQuery(place, timeSpans(dateNames), categoryIds)
	return web.database.Table("event").Count(query)
}

func (web *WebServer) countOrganizers() (int, error) {

	return web.database.Table("user").Count(nil)
}

func (web *WebServer) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		eventCnt, err := web.countEvents("Berlin", []int{0}, []string{"aktuell"})
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := web.countOrganizers()
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return web.view("start.tpl", w, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "categories": CategoryOrder, "categoryIds": CategoryMap})
	}()

	web.handle(w, result)
}

func (web *WebServer) approvePage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		var err error = nil

		if bson.IsObjectIdHex(r.Param("id")) {
			var user User
			userId := bson.ObjectIdHex(r.Param("id"))

			err = web.database.Table("user").LoadById(userId, &user)
			if err == nil {
				user.Approved = true
				if _, err = web.database.Table("user").UpsertById(userId, &user); err != nil {
					return &webResult{Status: http.StatusInternalServerError, Error: err}
				}
			}
		}

		return web.view("approve.tpl", w, bson.M{"approved": err == nil})
	}()

	web.handle(w, result)
}

func (web *WebServer) eventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	radius, err := strconv.Atoi(r.Param("radius"))
	if err != nil {
		radius = 0
	}

	place := r.Param("place")
	dateNames := strings.Split(r.Param("dates"), ",")
	categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))

	result := func() *webResult {

		eventCnt, err := web.countEvents(place, categoryIds, dateNames)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		organizerCnt, err := web.countOrganizers()
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		var result EventSearchResult
		query := buildQuery(place, timeSpans(dateNames), categoryIds)
		err = web.database.Table("event").Search(query, 0, 10, &result, "start")
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return web.view("events.tpl", w, bson.M{"eventCnt": eventCnt, "organizerCnt": organizerCnt, "events": result.Events, "place": place, "radius": radius, "dates": dateNames, "categoryIds": categoryIds, "categories": CategoryOrder, "categoryMap": CategoryMap})
	}()

	web.handle(w, result)
}

func (web *WebServer) checkSession(r *Request) (*User, error) {

	sessionId, err := r.ReadSessionId()
	if err != nil {
		return nil, err
	}

	user, err := web.database.LoadUserBySessionId(sessionId)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (web *WebServer) adminPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession((&Request{r}))
		if err != nil {
			return resultBadRequest
		}

		var events []Event
		err = web.database.Table("event").Find(bson.M{"organizerid": user.Id}, &events, "-start")
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return web.view("admin.tpl", w, bson.M{"user": user, "events": events})
	}()

	web.handle(w, result)
}

func (web *WebServer) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession((&Request{r}))
		if err != nil {
			return resultBadRequest
		}

		return web.view("profile.tpl", w, bson.M{"user": user})
	}()

	web.handle(w, result)
}

func (web *WebServer) passwordPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession((&Request{r}))
		if err != nil {
			return resultBadRequest
		}

		return web.view("password.tpl", w, bson.M{"user": user})
	}()

	web.handle(w, result)
}

func (web *WebServer) eventPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession((&Request{r}))
		if err != nil {
			return resultBadRequest
		}

		var data Event
		event := &data
		if bson.IsObjectIdHex(r.Param("id")) {

			err = web.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
			if err != nil {
				return &webResult{Status: http.StatusNotFound, Error: err}
			}

			if event.OrganizerId != user.Id {
				return resultBadRequest
			}
		}

		return web.view("event.tpl", w, bson.M{"user": user, "event": event, "categories": CategoryOrder, "categoryIds": CategoryMap})
	}()

	web.handle(w, result)
}

func isEmpty(s string) bool {

	return len(strings.TrimSpace(s)) == 0
}

func (web *WebServer) niceUrl(s string) string {

	return strings.ToLower(strings.Trim(web.niceExpr.ReplaceAllString(s, "-"), "-"))
}

func (web *WebServer) searchHandler(w traffic.ResponseWriter, r *traffic.Request) {

	path := r.PostFormValue("search")
	if path == "organizers" {
		path = "/veranstalter/"
	} else {
		path = "/veranstaltungen/"
	}

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
	categoryNames := make([]string, len(categoryIds))
	for i, id := range categoryIds {
		categoryNames[i] = CategoryIdMap[id]
	}

	dateIds := str2Int(r.Form["date"])
	if len(dateIds) == 0 {
		dateIds = append(dateIds, 0)
	}

	dateNames := make([]string, len(dateIds))
	for i, id := range dateIds {
		dateNames[i] = DateIdMap[id]
	}

	w.Header().Set("Location", path+place+"/"+strings.Join(dateNames, ",")+"/"+strings.Join(int2Str(categoryIds), ",")+"/"+strconv.Itoa(radius)+"/"+strings.Join(categoryNames, ","))
	w.WriteHeader(http.StatusFound)
}

func (web *WebServer) uploadHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		file, _, err := r.FormFile("file")
		if err != nil {
			return resultBadRequest
		}

		filename := uuid.New() + ".jpg"
		resp, err := http.Post(web.imgServer+"/"+filename, "image/jpeg", file)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		if resp.StatusCode == http.StatusOK {
			return &webResult{Status: http.StatusOK, JSON: filename}
		} else {
			return &webResult{Status: resp.StatusCode}
		}
	}()

	web.handle(w, result)
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

func (web *WebServer) sendEmail(to, subject, body string) error {

	tpl, err := web.tpls.Find("email.tpl")
	if err != nil {
		return err
	}

	return SendEmail(web.emailAccount, tpl, to, subject, body)
}

func (web *WebServer) registerHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := (&Request{r}).ReadUser()
		if err != nil {
			return resultBadRequest
		}

		err = validateUser(web.database, user)
		if err != nil {
			return resultConflict
		}

		user.Id = bson.NewObjectId()
		_, err = web.database.Table("user").UpsertById(user.Id, user)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = web.sendEmail(user.Email, register_subject, fmt.Sprintf(register_message, user.Addr.Name, user.Id.Hex()))
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		id, err := web.database.CreateSession(user.Id)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &webResult{Status: http.StatusCreated, JSON: id}
	}()

	web.handle(w, result)
}

func (web *WebServer) profileHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		request := &Request{r}

		user, err := web.checkSession(request)
		if err != nil {
			return resultBadRequest
		}

		data, err := request.ReadUser()
		if err != nil {
			return resultBadRequest
		}

		user.Addr.Name = data.Addr.Name
		user.Image = data.Image
		user.Descr = data.Descr
		user.Web = data.Web
		user.Addr.Street = data.Addr.Street
		user.Addr.Pcode = data.Addr.Pcode
		user.Addr.City = data.Addr.City

		_, err = web.database.Table("user").UpsertById(user.Id, user)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	web.handle(w, result)
}

func (web *WebServer) passwordHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		request := &Request{r}

		user, err := web.checkSession(request)
		if err != nil {
			return resultBadRequest
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
			err := web.sendEmail(user.Email, password_subject, fmt.Sprintf(password_message, user.Addr.Name, user.Id.Hex()))
			if err != nil {
				return &webResult{Status: http.StatusInternalServerError, Error: err}
			}
		}

		_, err = web.database.Table("user").UpsertById(user.Id, user)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	web.handle(w, result)
}

func (web *WebServer) unregisterHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession(&Request{r})
		if err != nil {
			return resultBadRequest
		}

		err = web.database.Table("event").Delete(bson.M{"organizerid": user.Id})
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		err = web.database.Table("user").DeleteById(user.Id)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	web.handle(w, result)
	web.logoutHandler(w, r)
}

func (web *WebServer) eventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		request := &Request{r}

		user, err := web.checkSession(request)
		if err != nil {
			return resultBadRequest
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
			web.database.Table("event").LoadById(event.Id, oldEvent)
			if oldEvent.OrganizerId != user.Id {
				return resultBadRequest
			}
		} else {
			created = true
			event.Id = bson.NewObjectId()
		}
		event.OrganizerId = user.Id

		_, err = web.database.Table("event").UpsertById(event.Id, event)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		if created {
			return resultCreated
		} else {
			return resultOK
		}
	}()

	web.handle(w, result)
}

func (web *WebServer) locationHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		locations := web.locations.Autocomplete(r.Param("location"))
		return &webResult{Status: http.StatusOK, JSON: locations}
	}()

	web.handle(w, result)
}

func (web *WebServer) eventCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		categoryIds := str2Int(strings.Split(r.Param("categoryIds"), ","))
		
		dateIds := str2Int(strings.Split(r.Param("dateIds"), ","))
		dates := make([]string, len(dateIds))
		for i, dateId := range dateIds {
			dates[i] = DateIdMap[dateId]
		}

		cnt, err := web.countEvents(r.Param("place"), categoryIds, dates)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &webResult{Status: http.StatusOK, JSON: cnt}
	}()

	web.handle(w, result)
}

func (web *WebServer) organizerCountHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		cnt, err := web.countOrganizers()
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &webResult{Status: http.StatusOK, JSON: cnt}
	}()

	web.handle(w, result)
}

func (web *WebServer) loginHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		form, err := (&Request{r}).ReadEmailAndPwd()
		if err != nil {
			return resultBadRequest
		}

		user, err := web.database.LoadUserByEmailAndPassword(form.Email, form.Pwd)
		if err != nil {
			return resultNotFound
		}

		id, err := web.database.CreateSession(user.Id)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return &webResult{Status: http.StatusCreated, JSON: id}
	}()

	web.handle(w, result)
}

func (web *WebServer) logoutHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		sessionId, err := (&Request{r}).ReadSessionId()
		if err != nil {
			return resultBadRequest
		}

		err = web.database.RemoveSession(sessionId)
		if err != nil {
			return resultNotFound
		}

		return resultOK
	}()

	web.handle(w, result)
}

func (web *WebServer) deleteEventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {

		user, err := web.checkSession(&Request{r})
		if err != nil || !bson.IsObjectIdHex(r.Param("id")) {
			return resultBadRequest
		}

		var data Event
		event := &data
		err = web.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
		if err != nil {
			return resultOK
		}

		if event.OrganizerId != user.Id {
			return resultBadRequest
		}

		err = web.database.Table("event").DeleteById(event.Id)
		if err != nil {
			return &webResult{Status: http.StatusInternalServerError, Error: err}
		}

		return resultOK
	}()

	web.handle(w, result)
}

func (web *WebServer) Start() {

	traffic.SetHost(web.host)
	traffic.SetPort(web.port)
	router := traffic.New()

	router.Get("/", web.startPage)
	router.Get("/approve/:id", web.approvePage)
	router.Get("/veranstalter/verwaltung", web.adminPage)
	router.Get("/veranstalter/verwaltung/kennwort", web.passwordPage)
	router.Get("/veranstalter/verwaltung/profil", web.profilePage)
	router.Get("/veranstalter/verwaltung/veranstaltung", web.eventPage)
	router.Get("/veranstalter/verwaltung/veranstaltung/:id", web.eventPage)
	router.Get("/veranstaltungen/:place/:dates/:categoryIds/:radius/:categories", web.eventsPage)

	router.Post("/suche", web.searchHandler)
	router.Post("/upload", web.uploadHandler)
	router.Post("/register", web.registerHandler)
	router.Post("/password", web.passwordHandler)
	router.Post("/profile", web.profileHandler)
	router.Post("/unregister", web.unregisterHandler)
	router.Post("/event", web.eventHandler)
	router.Get("/location/:location", web.locationHandler)
	router.Get("/eventcount/:place/:dateIds/:categoryIds", web.eventCountHandler)
	router.Get("/organizercount/:place/:categoryIds", web.organizerCountHandler)
	router.Post("/login", web.loginHandler)
	router.Post("/logout", web.logoutHandler)

	router.Delete("/event/:id", web.deleteEventHandler)

	router.Run()
}
