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
)

const (
	register_subject = "Deine Registrierung bei mitmachrepublik.de"
	register_message = "Liebe/r Organisator/in von %s,\r\n\r\nvielen Dank für die Registrierung bei der MitmachRepublik. Bitte bestätige Deine Registrierung, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://dev.mitmachrepublik.de/approve/%s\r\n\r\nDas Team der MitmachRepublik"
	password_subject = "Deine neue E-Mail-Adresse bei mitmachrepublik.de"
	password_message = "Liebe/r Organisator/in von %s,\r\n\r\nbitte bestätige Deine neue E-Mail-Adresse, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://dev.mitmachrepublik.de/approve/%s\r\n\r\nDas Team der MitmachRepublik"
)

func NewWebServer(host string, port int, tplDir, imgServer, mongoUrl, dbName string) (*WebServer, error) {

	database, err := NewMongoDb(mongoUrl, dbName)
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

func (web *WebServer) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	err := web.tpls.Execute("start.tpl", w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) approvePage(w traffic.ResponseWriter, r *traffic.Request) {

	var err error = nil

	if !bson.IsObjectIdHex(r.Param("id")) {
		err = errors.New("Failed to read id.")
	} else {
		var user User
		userId := bson.ObjectIdHex(r.Param("id"))

		err = web.database.Table("user").LoadById(userId, &user)
		if err == nil {
			user.Approved = true
			_, err = web.database.Table("user").UpsertById(userId, &user)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}
	}

	err = web.tpls.Execute("approve.tpl", w, bson.M{"approved": err == nil})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) eventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	err := web.tpls.Execute("events.tpl", w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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

	user, err := web.checkSession((&Request{r}))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var events []Event
	web.database.Table("event").Find(bson.M{"organizerid": user.Id}, &events)

	err = web.tpls.Execute("admin.tpl", w, bson.M{"user": user, "events": events})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := web.checkSession((&Request{r}))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.tpls.Execute("profile.tpl", w, bson.M{"user": user})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) passwordPage(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := web.checkSession((&Request{r}))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.tpls.Execute("password.tpl", w, bson.M{"user": user})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) eventPage(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := web.checkSession((&Request{r}))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data Event
	event := &data
	if bson.IsObjectIdHex(r.Param("id")) {

		err = web.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if event.OrganizerId != user.Id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}

	err = web.tpls.Execute("event.tpl", w, bson.M{"user": user, "event": event, "categories": CategoryOrder, "categoryIds": CategoryMap})
	if err != nil {
		fmt.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError)
	}
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

	file, _, err := r.FormFile("file")
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	filename := uuid.New() + ".jpg"
	resp, err := http.Post(web.imgServer+"/"+filename, "image/jpeg", file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if resp.StatusCode == http.StatusOK {
		w.WriteJSON(filename)
	} else {
		w.WriteHeader(resp.StatusCode)
	}
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

	user, err := (&Request{r}).ReadUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateUser(web.database, user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	user.Id = bson.NewObjectId()
	_, err = web.database.Table("user").UpsertById(user.Id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = web.sendEmail(user.Email, register_subject, fmt.Sprintf(register_message, user.Addr.Name, user.Id.Hex()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := web.database.CreateSession(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJSON(id)
}

func (web *WebServer) profileHandler(w traffic.ResponseWriter, r *traffic.Request) {

	request := &Request{r}

	user, err := web.checkSession(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := request.ReadUser()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
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
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) passwordHandler(w traffic.ResponseWriter, r *traffic.Request) {

	request := &Request{r}

	user, err := web.checkSession(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := request.ReadEmailAndPwd()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !isEmpty(data.Pwd) {
		user.Pwd = data.Pwd
	}

	if !isEmpty(data.Email) && data.Email != user.Email {
		user.Email = data.Email
		user.Approved = false
		err := web.sendEmail(user.Email, password_subject, fmt.Sprintf(password_message, user.Addr.Name, user.Id.Hex()))
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	_, err = web.database.Table("user").UpsertById(user.Id, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) unregisterHandler(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := web.checkSession(&Request{r})
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.database.Table("event").Delete(bson.M{"organizerid": user.Id})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = web.database.Table("user").DeleteById(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	web.logoutHandler(w, r)
}

func (web *WebServer) eventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	request := &Request{r}

	user, err := web.checkSession(request)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data Event
	event := &data
	err = request.ReadJson(event)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	created := false
	if event.Id.Valid() {
		var oldData Event
		oldEvent := &oldData
		web.database.Table("event").LoadById(event.Id, oldEvent)
		if oldEvent.OrganizerId != user.Id {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	} else {
		created = true
		event.Id = bson.NewObjectId()
	}
	event.OrganizerId = user.Id

	_, err = web.database.Table("event").UpsertById(event.Id, event)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if created {
		w.WriteHeader(http.StatusCreated)
	}
}

func (web *WebServer) loginHandler(w traffic.ResponseWriter, r *traffic.Request) {

	form, err := (&Request{r}).ReadEmailAndPwd()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := web.database.LoadUserByEmailAndPassword(form.Email, form.Pwd)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := web.database.CreateSession(user.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJSON(id)
}

func (web *WebServer) logoutHandler(w traffic.ResponseWriter, r *traffic.Request) {

	sessionId, err := (&Request{r}).ReadSessionId()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.database.RemoveSession(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (web *WebServer) deleteEventHandler(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := web.checkSession(&Request{r})
	if err != nil || !bson.IsObjectIdHex(r.Param("id")) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var data Event
	event := &data
	err = web.database.Table("event").LoadById(bson.ObjectIdHex(r.Param("id")), event)
	if err != nil {
		return /* http.StatusOK */
	}

	if event.OrganizerId != user.Id {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.database.Table("event").DeleteById(event.Id)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
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
