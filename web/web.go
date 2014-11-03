package mmr

import (
	"code.google.com/p/go-uuid/uuid"
	"encoding/json"
	"errors"
	"github.com/pilu/traffic"
	"html/template"
	"io/ioutil"
	"labix.org/v2/mgo/bson"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

type (
	WebServer struct {
		host      string
		port      int
		tplDir    string
		imgServer string
		database  Database
	}

	emailAndPwd struct {
		Email string
		Pwd   string
	}
)

func NewWebServer(host string, port int, tplDir, imgServer, mongoUrl, dbName string) (*WebServer, error) {

	database, err := NewMongoDb(mongoUrl, dbName)
	if err != nil {
		return nil, err
	}

	return &WebServer{host, port, tplDir, imgServer, database}, nil
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

func execute(template *template.Template, w traffic.ResponseWriter, data interface{}) {

	if template == nil {
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		template.Execute(w, data)
	}
}

func (web *WebServer) show(w traffic.ResponseWriter, view string, data interface{}) {

	templates, err := template.ParseGlob(web.tplDir + string(os.PathSeparator) + "*.tpl")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.WriteText(err.Error())
	} else {
		execute(templates.Lookup(view+".tpl"), w, data)
	}
}

func (web *WebServer) startPage(w traffic.ResponseWriter, r *traffic.Request) {

	web.show(w, "start", nil)
}

func (web *WebServer) eventsPage(w traffic.ResponseWriter, r *traffic.Request) {

	web.show(w, "events", nil)
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
	
	web.show(w, "admin", user)
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
	
	web.show(w, "profile", user)
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
	} else {
		w.WriteHeader(http.StatusCreated)
	}
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
