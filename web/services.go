package mmr

import (
	"encoding/json"
	"labix.org/v2/mgo/bson"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type (
	Service interface {
		Start()
		Stop()
	}

	BasicService struct {
		state int
		hour  int
		timer *time.Timer
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
		events *EventService
	}

	SendAlertsService struct {
		BasicService
		alerts *AlertService
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

func NewSessionService(hour int, database Database) Service {

	return &SessionService{NewBasicService(hour), database}
}

func (service *SessionService) Start() {

	service.start(service.serve)
}

func (service *SessionService) serve() {

	service.database.RemoveOldSessions(time.Duration(24) * time.Hour)
}

func NewUnusedImgService(hour int, database Database, imgServer string) Service {

	return &UnusedImgService{NewBasicService(hour), database, imgServer}
}

func (service *UnusedImgService) Start() {

	service.start(service.serve)
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

func (service *UnusedImgService) serve() {

	images, err := listImages(service.imgServer, 3600*24)
	if err != nil {
		return
	}

	var eventImages []string
	err = service.database.Table("event").Distinct(nil, "image", &eventImages)
	if err != nil {
		return
	}

	var userImages []string
	err = service.database.Table("user").Distinct(nil, "image", &userImages)
	if err != nil {
		return
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
			return
		}
		http.DefaultClient.Do(req)
	}
}

func NewUpdateRecurrencesService(hour int, events *EventService) Service {

	return &UpdateRecurrencesService{NewBasicService(hour), events}
}

func (service *UpdateRecurrencesService) Start() {

	service.start(service.serve)
}

func (service *UpdateRecurrencesService) serve() {

	service.events.UpdateRecurrences()
}

func NewSendAlertsService(hour int, alerts *AlertService) Service {

	return &SendAlertsService{NewBasicService(hour), alerts}
}

func (service *SendAlertsService) Start() {

	service.start(service.serve)
}

func (service *SendAlertsService) serve() {

}

func NewSpawnEventsService(hour int, database Database, events *EventService, imgServer string) Service {

	return &SpawnEventsService{NewBasicService(hour), database, events, imgServer}
}

func (service *SpawnEventsService) Start() {

	rand.Seed(time.Now().Unix())
	service.start(service.serve)
}

func (service *SpawnEventsService) serve() {

	titles := []string{"Volleyballtunier", "Wir haben es satt!", "Chor Open Stage Open Air", "Kinderbastelgruppe", "Jüdische Kulturtage", "Fit, Fun, Family im FEZ"}
	locations := []string{"Sportzentrum", "Brandenburger Tor", "Heiligengeistkirche", "Kindercafé", "Gemeindezentrum", "FEZ"}

	images, err := listImages(service.imgServer, 0)
	if err != nil {
		return
	}

	organizer, err := service.database.LoadUserByEmailAndPassword("leonhard.holz@web.de", "julius21")
	if err != nil {
		return
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
	service.events.Store(event, organizer.Approved)
}

func NewBasicService(hour int) BasicService {

	return BasicService{idle, hour, nil}
}

func timerDuration(hour int) time.Duration {

	now := time.Now()
	if now.Hour() >= hour {
		now.Add(time.Duration(24) * time.Hour)
	}
	start := time.Date(now.Year(), now.Month(), now.Day(), hour, 0, 0, 0, time.Local)
	return start.Sub(now)
}

func (service *BasicService) start(serve func()) {

	if service.state == idle {
		service.timer = time.NewTimer(timerDuration(service.hour))
		service.state = running
	}

	go func() {
		for {
			<-service.timer.C
			if service.state == running {
				serve()
				service.timer.Reset(timerDuration(service.hour))
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
