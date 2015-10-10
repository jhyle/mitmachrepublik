package mmr

import (
	"encoding/json"
	"github.com/pilu/traffic"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type (
	Service interface {
		Start()
		Run() error
		Stop()
	}

	BasicService struct {
		name  string
		state int
		hour  int
		timer *time.Timer
		email *EmailAccount
	}

	SessionService struct {
		BasicService
		database Database
	}

	UnusedImgService struct {
		BasicService
		database  Database
		imgServer string
	}

	UpdateRecurrencesService struct {
		BasicService
		events   *EventService
		account  *EmailAccount
		hostname string
	}

	SendAlertsService struct {
		BasicService
		hostname string
		from     *EmailAccount
		alerts   *AlertService
	}

	SpawnEventsService struct {
		BasicService
		database  Database
		events    *EventService
		imgServer string
	}
)

const (
	idle int = iota
	running
	stopping
)

var (
	alertSubject map[int]string = map[int]string{
		FromNow:     "Alle Veranstaltungen",
		Today:       "heute",
		Tomorrow:    "morgen",
		ThisWeek:    "diese Woche",
		NextWeekend: "am Wochenende",
		NextWeek:    "nächste Woche",
		TwoWeeks:    "in den nächsten 14 Tagen",
	}
)

func NewSessionService(hour int, email *EmailAccount, database Database) Service {

	return &SessionService{NewBasicService("SessionService", hour, email), database}
}

func (service *SessionService) Start() {

	service.start(service.Run)
}

func (service *SessionService) Run() error {

	return service.database.RemoveOldSessions(time.Duration(24) * time.Hour)
}

func NewUnusedImgService(hour int, email *EmailAccount, database Database, imgServer string) Service {

	return &UnusedImgService{NewBasicService("UnusedImagesService", hour, email), database, imgServer}
}

func (service *UnusedImgService) Start() {

	service.start(service.Run)
}

func listImages(imgServer string, age int) ([]string, error) {

	resp, err := http.Get(imgServer + "/?age=" + strconv.Itoa(age))
	if err != nil {
		return nil, err
	}

	var images []string
	err = json.NewDecoder(resp.Body).Decode(&images)
	if err != nil {
		return nil, err
	}

	return images, nil
}

func (service *UnusedImgService) Run() error {

	images, err := listImages(service.imgServer, 3600*24)
	if err != nil {
		return err
	}

	var eventImages []string
	err = service.database.Table("event").Distinct(nil, "image", &eventImages)
	if err != nil {
		return err
	}

	var userImages []string
	err = service.database.Table("user").Distinct(nil, "image", &userImages)
	if err != nil {
		return err
	}

	unusedImages := make(map[string]string)
	for _, image := range images {
		unusedImages[image] = image
	}

	for _, image := range eventImages {
		delete(unusedImages, image)
	}

	for _, image := range userImages {
		delete(unusedImages, image)
	}

	for image := range unusedImages {

		req, err := http.NewRequest("DELETE", service.imgServer+"/"+image, nil)
		if err != nil {
			return err
		}
		_, err = http.DefaultClient.Do(req)
		if err != nil {
			return err
		}
	}

	return nil
}

func NewUpdateRecurrencesService(hour int, email *EmailAccount, events *EventService, account *EmailAccount, hostname string) Service {

	return &UpdateRecurrencesService{NewBasicService("UpdateRecurrencesService", hour, email), events, account, hostname}
}

func (service *UpdateRecurrencesService) Start() {

	service.start(service.Run)
}

func (service *UpdateRecurrencesService) Run() error {

	dates, err := service.events.UpdateRecurrences()
	if err != nil {
		return err
	}

	if dates != nil {
		message := ""
		for _, dateId := range dates {
			date, err := service.events.LoadDate(dateId)
			if err == nil {
				message += "http://" + service.hostname + date.Url() + "\n"
			}
		}
		return SendEmail(service.account, service.account.From, nil, "Generierte Veranstaltungen", "text/plain", message)
	} else {
		return nil
	}
}

func NewSendAlertsService(hour int, email *EmailAccount, hostname string, from *EmailAccount, alerts *AlertService) *SendAlertsService {

	return &SendAlertsService{NewBasicService("SendAlertsService", hour, email), hostname, from, alerts}
}

func (service *SendAlertsService) Start() {

	service.start(service.Run)
}

func getNewsletter(hostname string, alert Alert) (string, error) {

	url := "/newsletter/veranstaltungen/" + eventSearchUrlWithQuery(alert.Place, alert.Targets, alert.Categories, alert.Dates, alert.Radius, alert.Query)
	resp, err := http.Get("http://" + hostname + url[0:strings.LastIndex(url, "/")+1] + alert.Id.Hex())
	if err != nil || resp.StatusCode != 200 {
		return "", err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	return string(bytes), err
}

func getNewsletterSubject(place string, dates []int) string {

	var subject string
	if len(dates) == 0 {
		subject = alertSubject[FromNow]
	} else if inArray(dates, TwoWeeks) {
		subject = "Veranstaltungen " + alertSubject[TwoWeeks]
	} else {
		i := 0
		names := make([]string, len(dates))
		for _, orderedDate := range DateOrder {
			if orderedDate != FromNow && orderedDate != TwoWeeks {
				for _, date := range dates {
					if orderedDate == date {
						names[i] = alertSubject[date]
						i++
						break
					}
				}
			}
		}
		subject = "Veranstaltungen " + strConcat(names)
	}

	if len(place) > 0 {
		subject += " in " + place
	}

	return subject
}

func (service *SendAlertsService) Run() error {

	alerts, err := service.alerts.FindAlerts(time.Now().Weekday())
	if err != nil {
		return err
	}

	for _, alert := range alerts {
		newsletter, err := getNewsletter(service.hostname, alert)
		if err != nil {
			return err
		} else if !isEmpty(newsletter) {
			err = SendEmail(service.from, &EmailAddress{alert.Name, alert.Email}, service.from.From, getNewsletterSubject(alert.Place, alert.Dates), "text/html", newsletter)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func NewSpawnEventsService(hour int, email *EmailAccount, database Database, events *EventService, imgServer string) Service {

	return &SpawnEventsService{NewBasicService("SpawnEventsService", hour, email), database, events, imgServer}
}

func (service *SpawnEventsService) Start() {

	rand.Seed(time.Now().Unix())
	service.start(service.Run)
}

func (service *SpawnEventsService) Run() error {

	titles := []string{"Volleyballtunier", "Wir haben es satt!", "Chor Open Stage Open Air", "Kinderbastelgruppe", "Jüdische Kulturtage", "Fit, Fun, Family im FEZ"}
	locations := []string{"Sportzentrum", "Brandenburger Tor", "Heiligengeistkirche", "Kindercafé", "Gemeindezentrum", "FEZ"}

	images, err := listImages(service.imgServer, 0)
	if err != nil {
		return err
	}

	organizer, err := service.database.LoadUserByEmailAndPassword("leonhard.holz@web.de", "julius21")
	if err != nil {
		return err
	}

	districts := make([]string, 0, len(PostcodeMap))
	for district := range PostcodeMap {
		districts = append(districts, district)
	}

	event := new(Event)
	event.Id = bson.NewObjectId()
	event.OrganizerId = organizer.Id
	event.Title = titles[rand.Intn(len(titles))]
	event.Image = images[rand.Intn(len(images))]
	event.Descr = "Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet. Lorem ipsum dolor sit amet, consetetur sadipscing elitr, sed diam nonumy eirmod tempor invidunt ut labore et dolore magna aliquyam erat, sed diam voluptua. At vero eos et accusam et justo duo dolores et ea rebum. Stet clita kasd gubergren, no sea takimata sanctus est Lorem ipsum dolor sit amet."
	event.Web = "http://www.facebook.com"
	event.Start = time.Now().Add(time.Duration(rand.Intn(720)) * time.Hour)
	event.Categories = make([]int, 0)
	for i := 0; i < rand.Intn(len(CategoryOrder)); i++ {
		event.Categories = append(event.Categories, CategoryMap[CategoryOrder[rand.Intn(len(CategoryOrder))]])
	}
	event.Addr.Name = locations[rand.Intn(len(locations))]
	event.Addr.Street = "Baker Street 221"
	district := districts[rand.Intn(len(districts))]
	event.Addr.Pcode = PostcodeMap[district][rand.Intn(len(PostcodeMap[district]))]
	event.Addr.City = "Berlin"
	return service.events.Store(event, organizer.Approved)
}

func NewBasicService(name string, hour int, email *EmailAccount) BasicService {

	return BasicService{name, idle, hour, nil, email}
}

func timerDuration(hour int) time.Duration {

	day := time.Now()
	if day.Hour() >= hour {
		day = day.AddDate(0, 0, 1)
	}
	start := time.Date(day.Year(), day.Month(), day.Day(), hour, 0, 0, 0, time.Local)
	return start.Sub(time.Now())
}

func (service *BasicService) start(serve func() error) {

	if service.state == idle {
		duration := timerDuration(service.hour)
		traffic.Logger().Printf("Starting timer for service %s, hour %v with duration %v.\n", service.name, service.hour, duration)
		service.timer = time.NewTimer(duration)
		service.state = running
	}

	go func() {
		for {
			<-service.timer.C
			traffic.Logger().Printf("Fired timer for service %s, hour %v, running = %v.\n", service.name, service.hour, running)
			if service.state == running {
				err := serve()
				if err != nil {
					traffic.Logger().Printf(service.name + ": " + err.Error())
					SendEmail(service.email, service.email.From, nil, "Fehlermeldung von " + service.name, "text/plain", err.Error())
				}
				duration := timerDuration(service.hour)
				traffic.Logger().Printf("Reseting timer for service %s, hour %v with duration %v.\n", service.name, service.hour, duration)
				service.timer.Reset(duration)
			} else {
				service.timer = nil
				service.state = idle
				break
			}
		}
	}()
}

func (service *BasicService) Stop() {

	if service.state == running {
		service.state = stopping
	}
}
