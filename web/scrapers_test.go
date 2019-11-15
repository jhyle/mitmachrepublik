package mmr_test

import (
	"testing"

	"github.com/jhyle/mitmachrepublik/web"
)

const (
	testMongoUrl = "localhost"
	testDatabase = "mitmachrepublik"
	hostname     = "https://www.mitmachrepublik.de"
)

func TestScraperService(t *testing.T) {

	database, err := mmr.NewMongoDb(testMongoUrl, testDatabase)
	if err != nil {
		t.Fatal(err)
	}

	users, err := mmr.NewUserService(database, "user")
	if err != nil {
		t.Fatal(err)
	}

	events, err := mmr.NewEventService(database, "event", "/tmp/mmr/index")
	if err != nil {
		t.Fatal(err)
	}

	admin, err := users.LoadByEmail(mmr.ADMIN_EMAIL)
	if err != nil {
		t.Fatal(err)
	}

	err = mmr.NewScrapersService(0, nil, events, admin.Id).Run()
	if err != nil {
		t.Fatal(err)
	}
}
