package main

import (
	"twitterstreamprocessor/twitter"
	"twitterstreamprocessor/worker"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

const parallelism = 4

func main() {
	conf := twitter.GetConfig()
	ch := make(chan map[string]int)
	cnts := make(map[string]int)

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	for i := 0; i < parallelism; i++ {
		worker := worker.NewWorker(conf, ch, i, parallelism)
		go worker.Start()
	}

	maxcnt := 0
	currcnt := 0
	maxtag := ""

	go func() {
		<-sigs
		fmt.Printf("\nThe most popular hashtag was %s with %d occurrences.\n", maxtag, maxcnt)
		done <- true
	}()

	for {
		select {
		case msg := <-ch:
			for k,v := range msg {
				if cnt, ok := cnts[k]; ok {
					cnts[k] = cnt+v
					currcnt = cnt+v
				} else {
					cnts[k] = v
					currcnt = v
				}
				if currcnt > maxcnt {
					maxcnt = currcnt
					maxtag = k
				}
			}
		case <-done:
			fmt.Println("Done.")
			return
		}
	}
}
