package mmr_test

import (
	"flag"
	"testing"

	"github.com/jhyle/mitmachrepublik/web"

	"gopkg.in/mgo.v2/bson"
)

var (
	gibUser     *string = flag.String("gibUser", "", "Gratis in Berlin user")
	gibPassword *string = flag.String("gibPassword", "", "Gratis in Berlin password")

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
		GiBId:      "2022206",
	}
)

func TestGratisInBerlin(t *testing.T) {

	client, err := mmr.NewGibClient(hostname, *gibUser, *gibPassword)
	if err != nil {
		t.Fatal(err)
	}

	event.GiBId, err = client.PostEvent(&event)
	if err != nil {
		t.Fatal(err)
	}

	err = client.DeletePost(&event)
	if err != nil {
		t.Fatal(err)
	}
}
