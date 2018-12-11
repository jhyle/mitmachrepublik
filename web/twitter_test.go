package mmr_test

import (
	"flag"
	"testing"

	"github.com/jhyle/mitmachrepublik/web"
)

var (
	twitterApiSecret         *string = flag.String("twitterApiSecret", "", "Twitter Api Secret")
	twitterAccessTokenSecret *string = flag.String("twitterAccessTokenSecret", "", "Twitter Api Access Token Secret")
)

func TestTwitterPost(t *testing.T) {

	client := mmr.NewTwitterClient(hostname, *twitterApiSecret, *twitterAccessTokenSecret)

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
