package main

import (
	"log"
	"sync"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func getEmails(messages chan *imap.Message, wg *sync.WaitGroup, from, to int) {

	//log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.yandex.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Connected to server")

	// Don't forget to logout
	defer c.Logout()
	defer log.Println("Logging out")
	defer wg.Done()

	// Login
	if err := c.Login("stuffsom", "some_password1"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// Select INBOX
	_, err = c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("number of messages : ", mbox.Messages)

	seqset := new(imap.SeqSet)
	seqset.AddRange(uint32(from), uint32(to))

	log.Println("executing for seqset", seqset.String())
	err = c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	var MessageMap sync.Map
	var wg sync.WaitGroup
	NumEmailsRead := 0

	start := time.Now()
	defer func() {
		log.Printf("Read %v messages", NumEmailsRead)
		elapsed := time.Since(start)
		log.Printf("Time taken %s", elapsed)
	}()

	const NumEmails = 500
	for i := 1; i <= 2000; i = i + NumEmails {
		wg.Add(1)
		m := make(chan *imap.Message)
		go getEmails(m, &wg, i, i+(NumEmails-1))

		go func(m chan *imap.Message) {
			for msg := range m {
				MessageMap.Store(msg.Envelope.Subject, 1)
				NumEmailsRead++
			}
		}(m)
	}

	wg.Wait()
}
