package mmr

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/pilu/traffic"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"fmt"
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
	register_subject = "Bestätige Deine Registrierung"
	register_message = "Liebe/r Organisator/in von %s,\r\n\r\nvielen Dank für die Registrierung bei der MitmachRepublik. Bitte bestätige Deine Registrierung, in dem Du auf den folgenden Link klickst:\r\n\r\nhttp://dev.mitmachrepublik.de/approve/%s\r\n\r\nDas Team der MitmachRepublik"
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

func readBinary(r *traffic.Request) ([]byte, error) {

	return ioutil.ReadAll(r.Request.Body)
}

func readJson(r *traffic.Request, v interface{}) error {

	buffer, err := readBinary(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(buffer, v)
}

func readSessionId(r *traffic.Request) (bson.ObjectId, error) {

	cookie, err := r.Cookie("SESSIONID")
	if err != nil || !bson.IsObjectIdHex(cookie.Value) {
		return "", errors.New("Failed to read session id.")
	}

	return bson.ObjectIdHex(cookie.Value), nil
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

	err = web.tpls.Execute("approve.tpl", w, err == nil)
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

func (web *WebServer) adminPage(w traffic.ResponseWriter, r *traffic.Request) {

	sessionId, err := readSessionId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := web.database.LoadUserBySessionId(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = web.tpls.Execute("admin.tpl", w, user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) profilePage(w traffic.ResponseWriter, r *traffic.Request) {

	sessionId, err := readSessionId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := web.database.LoadUserBySessionId(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	err = web.tpls.Execute("profile.tpl", w, user)
	if err != nil {
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

func readUser(r *traffic.Request) (*User, error) {

	var data User
	user := &data
	err := readJson(r, user)
	return user, err
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

func (web *WebServer) registerHandler(w traffic.ResponseWriter, r *traffic.Request) {

	user, err := readUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = validateUser(web.database, user)
	if err != nil {
		w.WriteHeader(http.StatusConflict)
		return
	}

	user.SetId(bson.NewObjectId())
	_, err = web.database.Table("user").UpsertById(user.GetId(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	tpl, err := web.tpls.Find("email.tpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	err = SendEmail(web.emailAccount, tpl, user.Email, register_subject, fmt.Sprintf(register_message, user.Addr.Name, user.Id.Hex()))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	id, err := web.database.CreateSession(user.GetId())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJSON(id)
}

func (web *WebServer) profileHandler(w traffic.ResponseWriter, r *traffic.Request) {

	sessionId, err := readSessionId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, err := readUser(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := web.database.LoadUserBySessionId(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	user.Addr.Name = data.Addr.Name
	user.Image = data.Image
	user.Descr = data.Descr
	user.Web = data.Web
	user.Addr.Street = data.Addr.Street
	user.Addr.Pcode = data.Addr.Pcode
	user.Addr.City = data.Addr.City

	_, err = web.database.Table("user").UpsertById(user.GetId(), user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (web *WebServer) loginHandler(w traffic.ResponseWriter, r *traffic.Request) {

	var form emailAndPwd
	err := readJson(r, &form)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := web.database.LoadUserByEmailAndPassword(form.Email, form.Pwd)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	id, err := web.database.CreateSession(user.GetId())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.WriteJSON(id)
}

func (web *WebServer) logoutHandler(w traffic.ResponseWriter, r *traffic.Request) {

	sessionId, err := readSessionId(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = web.database.RemoveSession(sessionId)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
}

func (web *WebServer) Start() {

	traffic.SetHost(web.host)
	traffic.SetPort(web.port)
	router := traffic.New()

	router.Get("/", web.startPage)
	router.Get("/approve/:id", web.approvePage)
	router.Get("/veranstalter/verwaltung", web.adminPage)
	router.Get("/veranstalter/verwaltung/profil", web.profilePage)
	router.Get("/veranstaltungen/:place/:radius/:category", web.eventsPage)

	router.Post("/suche", web.searchHandler)
	router.Post("/upload", web.uploadHandler)
	router.Post("/register", web.registerHandler)
	router.Post("/profile", web.profileHandler)
	router.Post("/login", web.loginHandler)
	router.Post("/logout", web.logoutHandler)

	router.Run()
}
