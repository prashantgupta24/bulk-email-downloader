package main

import (
	"log"
	"sync"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
)

func main() {
	log.Println("Connecting to server...")

	// Connect to server
	c, err := client.DialTLS("imap.yandex.com:993", nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected")

	// Don't forget to logout
	defer c.Logout()
	defer log.Println("Logging out")

	// Login
	if err := c.Login("stuffsom", "some_password1"); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	// List mailboxes
	// mailboxes := make(chan *imap.MailboxInfo, 10)
	// //done := make(chan error, 1)

	// err = c.List("", "*", mailboxes)

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// // go func() {
	// // 	done <- c.List("", "*", mailboxes)
	// // }()

	// log.Println("Mailboxes:")
	// for m := range mailboxes {
	// 	log.Println("* " + m.Name)
	// }

	// if err := <-done; err != nil {
	// 	log.Fatal(err)
	// }

	// Select INBOX
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	//log.Println("Flags for INBOX:", mbox.Flags)
	log.Println("number of messages : ", mbox.Messages)

	// //fetch all messages
	// from := uint32(1)
	// //to := mbox.Messages
	// to := uint32(1)
	// seqset := new(imap.SeqSet)
	// seqset.AddRange(from, to)
	messages := make(chan *imap.Message)
	// errChan := make(chan error)
	MessageMap := make(map[string]int)
	var wg sync.WaitGroup

	type DetailedError struct {
		errVal error
		seqset string
	}

	start := time.Now()
	defer func() {
		log.Printf("Read %v messages", len(MessageMap))
		elapsed := time.Since(start)
		log.Printf("Time taken %s", elapsed)
	}()

	const NumEmailsPerWorker = 1000
	for i := uint32(1); i <= mbox.Messages; i = i + NumEmailsPerWorker {
		wg.Add(1)
		m := make(chan *imap.Message)
		e := make(chan DetailedError)
		seqset := new(imap.SeqSet)
		seqset.AddRange(i, i+(NumEmailsPerWorker-1))
		//seqset.AddRange(i, mbox.Messages)

		go func(seqset *imap.SeqSet, m chan *imap.Message, e chan DetailedError) {
			defer wg.Done()
			defer close(e)

			log.Println("executing for seqset", seqset.String())
			errVal := c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, m)
			e <- DetailedError{
				errVal: errVal,
				seqset: seqset.String(),
			}
		}(seqset, m, e)

		go func(m chan *imap.Message, e chan DetailedError) {
			for msg := range m {
				//log.Println("message is ", msg)
				messages <- msg
			}
			for DetailedError := range e {
				if DetailedError.errVal != nil {
					log.Println("fatal error for seqset", DetailedError.seqset)
					// log.Fatal(err)
				}
			}

		}(m, e)
	}
	// go func() {
	// 	errChan <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	// }()

	// go func() {
	// 	for msg := range messages {
	// 		MessageMap[msg.Envelope.Subject] = 1
	// 		//log.Println("* Message read : " + msg.Envelope.Subject)
	// 		//wg.Done()
	// 	}
	// }()

	wg.Wait()

	// if err := <-errChan; err != nil {
	// 	log.Fatal(err)
	// }

	// // Get the last 4 messages
	// from := uint32(1)
	// to := mbox.Messages
	// if mbox.Messages > 3 {
	// 	// We're using unsigned integers here, only substract if the result is > 0
	// 	from = mbox.Messages - 3
	// }
	// seqset := new(imap.SeqSet)
	// seqset.AddRange(from, to)

	// messages := make(chan *imap.Message, 10)
	// done := make(chan error, 1)
	// go func() {
	// 	done <- c.Fetch(seqset, []imap.FetchItem{imap.FetchEnvelope}, messages)
	// }()

	// log.Println("Last 4 messages:")
	// for msg := range messages {
	// 	log.Println("* " + msg.Envelope.Subject)
	// }

	// if err := <-done; err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Done!")
}
