package main

import (
	"log"
	"os"
	"sync"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

func getEmails(messages chan *imap.Message, wg *sync.WaitGroup, from, to int) {

	log.Println("Connecting to server...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Connect to server
	c, err := client.DialTLS(os.Getenv("IMAP_SERVER"), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Don't forget to logout
	defer c.Logout()
	defer log.Println("Logging out")

	// Login
	if err := c.Login(os.Getenv("EMAIL"), os.Getenv("PASSWORD")); err != nil {
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

	log.Println("Fetching emails from number", seqset.String())
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

	const EmailsPerBatch = 500
	const TotalEmails = 2000

	for i := 1; i <= TotalEmails; i = i + EmailsPerBatch {
		wg.Add(1)
		m := make(chan *imap.Message)

		go getEmails(m, &wg, i, i+(EmailsPerBatch-1))

		go func(m chan *imap.Message) {
			for msg := range m {
				MessageMap.Store(msg.Envelope.Subject, 1)
				NumEmailsRead++
			}
			wg.Done()
		}(m)
	}

	wg.Wait()
}
