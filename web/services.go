package mmr

import (
	"encoding/json"
	"net/http"
	"time"
	"strconv"
)

type (
	Service interface {
		Start()
		Stop()
	}

	BasicService struct {
		state    int
		interval int
		ticker   *time.Ticker
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
)

const (
	idle int = iota
	running
	stopping
)

func NewSessionService(interval int, database Database) Service {

	return &SessionService{BasicService{idle, interval, nil}, database}
}

func (service *SessionService) Start() {

	service.start(service.serve)
}

func (service *SessionService) serve() {

	service.database.RemoveOldSessions(time.Duration(24) * time.Hour)
}

func NewUnusedImgService(interval int, database Database, imgServer string) Service {

	return &UnusedImgService{BasicService{idle, interval, nil}, database, imgServer}
}

func (service *UnusedImgService) Start() {

	service.start(service.serve)
}

func (service *UnusedImgService) serve() {

	resp, err := http.Get(service.imgServer + "/?age=" + strconv.Itoa(3600 * 24))
	if err != nil {
		return
	}

	var images []string
	err = json.NewDecoder(resp.Body).Decode(&images)
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

		req, err := http.NewRequest("DELETE", service.imgServer + "/" + image, nil)
		if err != nil {
			return
		}

		resp, err = http.DefaultClient.Do(req)
		if err != nil {
			return
		}
	}
}

func (service *BasicService) start(serve func()) {

	if service.state == idle {
		service.ticker = time.NewTicker(time.Duration(service.interval) * time.Second)
		service.state = running
	}

	go func() {
		for _ = range service.ticker.C {
			if service.state == stopping {
				break
			}
			serve()
		}
		service.ticker = nil
		service.state = idle
	}()
}

func (service *BasicService) Stop() {

	if service.state == running {
		service.state = stopping
	}
}
