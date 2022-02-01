package main

import (
	"bufio"
	"bytes"
	"github.com/Rhymen/go-whatsapp"
	"github.com/skip2/go-qrcode"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

const text = "Hello! This is Hammed Lekan Adigun. I am contacting you with my new WhatsApp phone number. Would you find it in your kind heart to save it in your contact book? \n\nThank you!**Kindly ignore if you’ve saved it already."

func main() {
	b, err := ioutil.ReadFile("data.txt")
	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewScanner(bytes.NewBuffer(b))

	type Contact struct {
		Name string
		Phone string
	}
	contacts := make([]*Contact, 0)
	for reader.Scan() {
		ss := strings.Split(reader.Text(), ",")
		contacts = append(contacts, &Contact{Name:ss[0], Phone: ss[1]})
	}
	qrCodePath := "code.jpeg"
	conn, err := whatsapp.NewConn(30 * time.Second)
	if err != nil {
		log.Fatal(err)
	}
	qrChan := make(chan string)
	go func() {
		code := <-qrChan
		err := qrcode.WriteFile(code, qrcode.Medium, 256, qrCodePath)
		if err != nil {
			log.Fatalf("failed to write qrCode: %v", err)
		}
	}()

	sess, err := conn.Login(qrChan)
	if err != nil {
		log.Fatalf("failed to login: %v", err)
	}
	log.Printf("Login sess established: %v", sess)
	for _, nextContact := range contacts {
		textMessage := whatsapp.TextMessage{
			Info:        whatsapp.MessageInfo{
				RemoteJid: nextContact.Phone + "@s.whatsapp.net",
			},
			Text:        text,
		}
		time.Sleep(1 * time.Second)
		log.Printf("sending text to: %s->%s", nextContact.Name, nextContact.Phone)
		resp, err := conn.Send(textMessage)
		if err != nil {
			log.Printf("failed to send message to %s: %v", nextContact.Name, err)
		}
		log.Printf("Message sent to %s! Response=%s", nextContact.Name, resp)
	}
}
