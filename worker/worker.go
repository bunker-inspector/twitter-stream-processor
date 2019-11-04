package worker

import (
	"twitterstreamprocessor/twitter"
	t "github.com/dghubble/go-twitter/twitter"
	"strings"
	"fmt"
)

type Worker struct {
	Client *twitter.StreamClient
	Channel chan map[string]int
	Offset int
	Parallelism int
}

func NewWorker(conf* twitter.Config,
	ch chan map[string]int,
	offset int,
	parallelism int) *Worker {
	fmt.Println(offset)
	return &Worker{
		Client: twitter.NewClient(conf),
		Channel: ch,
		Offset: offset,
		Parallelism: parallelism}
}

func (w* Worker) Start() {
	stream, err := w.Client.Start()

	if err != nil {
		fmt.Println("Error.")
	}

	loopcnt := 0
	for rawMessage := range stream.Messages {
		loopcnt++
		loopcnt %= w.Parallelism
		if loopcnt != w.Offset {
			continue
		}

		switch message := rawMessage.(type) {
		case *t.Tweet:
			fmt.Println(message.Text)
			cnts := make(map[string]int)

			for _, word := range strings.Split(message.Text, " ") {
				if strings.HasPrefix(word, "#") {
					hashtag := word[1:]
					if cnt, ok := cnts[hashtag]; ok {
						cnts[hashtag] = cnt + 1
					} else {
						cnts[hashtag] = 1
					}
				}
			}
			w.Channel <- cnts
		default:
			continue
		}
	}
}
