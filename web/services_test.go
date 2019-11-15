package mmr_test

import (
	"testing"

	"github.com/jhyle/mitmachrepublik/web"
)

func TestRemoveObsoleteEventsService(t *testing.T) {

	database, err := mmr.NewMongoDb(testMongoUrl, testDatabase)
	if err != nil {
		t.Fatal(err)
	}

	err = mmr.NewObsoleteEventsService(0, nil, database).Run()
	if err != nil {
		t.Fatal(err)
	}
}
