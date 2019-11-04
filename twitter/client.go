package twitter

import (
	"github.com/BurntSushi/toml"
	"os"
	"fmt"
	//"net/http"
	t "github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
)

const BASE = "https://api.twitter.com"
const STREAM = "https://stream.twitter.com/1.1"

type StreamClient struct {
	c *t.Client
}

type Config struct {
	ConsumerKey    string `toml:consumerKey`
	ConsumerSecret string `toml:consumerSecret`
	AccessToken    string `toml:accessToken`
	AccessSecret   string `toml:accessSecret`
}

func GetConfig() *Config {
	var config Config
	path := fmt.Sprintf("%s/config/twitter.toml", os.Getenv("GOPATH"))
	if _, err := toml.DecodeFile(path, &config); err != nil {
		panic(err.Error())
	}
	return &config
}

func NewClient(c *Config) *StreamClient {
	conf := oauth1.NewConfig(c.ConsumerKey, c.ConsumerSecret)
	token := oauth1.NewToken(c.AccessToken, c.AccessSecret)
	httpClient := conf.Client(oauth1.NoContext, token)

	return &StreamClient{c: t.NewClient(httpClient)}
}

func (c *StreamClient) Start() (*t.Stream, error) {
	return c.c.Streams.Sample(&t.StreamSampleParams{})
}
