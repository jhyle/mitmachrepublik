package mmr

import (
	"time"

	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	"github.com/pkg/errors"
)

const (
	twitterApiKey      = "7L6j8DAByj9Hfv1tbAnccJtNz"
	twitterAccessToken = "3177504856-2nuAxUOxRyMtkqLOsRMlPEJ5AzxvIj2fXKa42HP"
)

type TwitterClient struct {
	hostname string
	client   *twitter.Client
}

func NewTwitterClient(hostname, apiSecret, accessTokenSecret string) *TwitterClient {

	config := oauth1.NewConfig(twitterApiKey, apiSecret)
	token := oauth1.NewToken(twitterAccessToken, accessTokenSecret)
	httpClient := config.Client(oauth1.NoContext, token)

	return &TwitterClient{hostname, twitter.NewClient(httpClient)}
}

func (twttr *TwitterClient) PostEvent(event *Event) (int64, error) {

	nextDate := event.NextDate(time.Now())
	message := event.Title + " am " + dateFormat(nextDate) + " um " + timeFormat(nextDate) + " in " + event.Addr.City + ": " + twttr.hostname + event.Url()

	tweet, _, err := twttr.client.Statuses.Update(message, nil)
	if err != nil {
		return 0, errors.Wrap(err, "error sending tweet")
	}

	return tweet.ID, nil
}

func (twttr *TwitterClient) DeletePost(id int64) error {

	_, _, err := twttr.client.Statuses.Destroy(id, nil)
	if err != nil {
		return errors.Wrap(err, "error deleting tweet")
	}

	return nil
}
