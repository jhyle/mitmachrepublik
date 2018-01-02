package mmr_test

import (
	"flag"
	"testing"
	"time"

	"github.com/jhyle/mmr/web"
	"gopkg.in/mgo.v2/bson"
)

const (
	hostname = "https://www.mitmachrepublik.de"
	appId    = "138725613479008"
)

var (
	event = mmr.Event{
		Id:          bson.ObjectIdHex("568d17f79f21890832055c3e"),
		OrganizerId: bson.ObjectIdHex("554da1609f2189538b000001"),
		Title:       "Wir haben es satt!",
		Image:       "8948e387-43d2-49ce-8916-9346ffb94871.jpg",
		Descr:       "Die Landwirtschaft steht am Scheideweg: Wird unser Essen zukünftig noch von Bäuerinnen und Bauern erzeugt oder von Agrarkonzernen, die auf Agrogentechnik und Tierfabriken setzen und zu Dumpingpreisen für den Weltmarkt produzieren?",
		Addr: mmr.Address{
			Name:   "Brandenburger Tor",
			Street: "Pariser Platz",
			Pcode:  "10117",
			City:   "Berlin",
		},
		Targets:    []int{2, 3, 4, 5, 1, 6},
		Categories: []int{7, 13, 14, 23, 15},
	}

	fbAppSecret *string = flag.String("fbAppSecret", "", "Facebook App Secret")
	fbUser      *string = flag.String("fbUser", "", "Facebook user")
	fbPassword  *string = flag.String("fbPassword", "", "Facebook password")
)

func TestFacebookPost(t *testing.T) {

	client, err := mmr.NewFacebookClient(hostname, appId, *fbAppSecret, *fbUser, *fbPassword)
	if err != nil {
		t.Fatal(err)
	}

	id, err := client.PostEvent(&event)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(id)

	err = client.DeletePost(id)
	if err != nil {
		t.Fatal(err)
	}
}

func init() {

	t, err := time.Parse(time.RFC3339Nano, "2018-01-20T12:00:00.000Z")
	if err != nil {
		panic(err)
	}
	event.Start = t
}
