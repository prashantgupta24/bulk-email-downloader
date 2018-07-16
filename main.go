package main

import (
	"io/ioutil"
	"log"
	"net/mail"
	"os"
	"sync"
	"time"

	imap "github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
)

//Email message constructed
type MessageStruct struct {
	From, Subject, Body string
}

func getEmails(MessageChan chan *imap.Message, from, to int, section *imap.BodySectionName) {

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
	err = c.Fetch(seqset, []imap.FetchItem{section.FetchItem()}, MessageChan)
	if err != nil {
		log.Println(err)
	}
}

func main() {

	var Messages []MessageStruct
	var wg sync.WaitGroup
	section := imap.BodySectionName{}

	start := time.Now()
	defer func() {
		log.Printf("Read %v messages", len(Messages))
		elapsed := time.Since(start)
		log.Printf("Time taken %s", elapsed)
	}()

	const EmailsPerBatch = 500
	const TotalEmails = 2000

	for i := 1; i <= TotalEmails; i = i + EmailsPerBatch {
		wg.Add(1)
		MessageChan := make(chan *imap.Message)

		go getEmails(MessageChan, i, i+(EmailsPerBatch-1), &section)

		go func() {
			for RawMessage := range MessageChan {
				r := RawMessage.GetBody(&section)
				if r == nil {
					log.Println("Server didn't return message body")
				} else {
					Message, err := mail.ReadMessage(r)
					if err != nil {
						log.Fatal(err)
					}

					header := Message.Header
					body, err := ioutil.ReadAll(Message.Body)
					if err != nil {
						log.Fatal(err)
					}
					Messages = append(Messages, MessageStruct{
						From:    header.Get("From"),
						Subject: header.Get("Subject"),
						Body:    string(body),
					})
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
