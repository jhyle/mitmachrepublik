package mmr

import (
	"code.google.com/p/go-uuid/uuid"
	"errors"
	"fmt"
	"github.com/pilu/traffic"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type (
	WebServer struct {
		host         string
		port         int
		tpls         *Templates
		imgServer    string
		database     Database
		emailAccount *EmailAccount
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
	return &WebServer{host, port, tpls, imgServer, database, emailAccount}, nil
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

func (web *WebServer) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	result := func() *webResult {
		return web.view("start.tpl", w, nil)
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

	result := func() *webResult {
		return web.view("events.tpl", w, nil)
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

func (web *WebServer) searchHandler(w traffic.ResponseWriter, r *traffic.Request) {

	path := r.PostFormValue("search")
	if path == "organizers" {
		path = "/veranstalter/"
	} else {
		path = "/veranstaltungen/"
	}

	place := strings.ToLower(r.PostFormValue("place"))
	if isEmpty(place) {
		place = "berlin"
	}

	radius, err := strconv.Atoi(r.PostFormValue("radius"))
	if err != nil {
		radius = 0
	}

	category := r.PostFormValue("category")

	w.Header().Set("Location", path+url.QueryEscape(place)+"/"+strconv.Itoa(radius)+"/"+url.QueryEscape(category))
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
	router.Get("/veranstaltungen/:place/:radius/:category", web.eventsPage)

	router.Post("/suche", web.searchHandler)
	router.Post("/upload", web.uploadHandler)
	router.Post("/register", web.registerHandler)
	router.Post("/password", web.passwordHandler)
	router.Post("/profile", web.profileHandler)
	router.Post("/unregister", web.unregisterHandler)
	router.Post("/event", web.eventHandler)
	router.Post("/login", web.loginHandler)
	router.Post("/logout", web.logoutHandler)

	router.Delete("/event/:id", web.deleteEventHandler)

	router.Run()
}
