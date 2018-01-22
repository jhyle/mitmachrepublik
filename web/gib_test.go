package mmr_test

import (
	"flag"
	"testing"

	"github.com/jhyle/mmr/web"
)

var (
	gibUser     *string = flag.String("gibUser", "", "Gratis in Berlin user")
	gibPassword *string = flag.String("gibPassword", "", "Gratis in Berlin password")
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
