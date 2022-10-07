# GoGmailnator
Simple emailnator.com wrapper made in Go.

- Create sessions, generate emails, receive mail
- Proxy Support

## Usage

```
package main

import (
	"github.com/ox-y/GoGmailnator"
	"fmt"
)

func main() {
	var sess GoGmailnator.Session

	proxy := "p.webshare.io:9999"

	// session will expire after a few hours
	err := sess.Init(&proxy)
	if err != nil {
		panic(err)
	}

	// calling sess.GenerateEmailAddress or sess.RetrieveMail with a dead session will cause an error
	isAlive, err := sess.IsAlive()
	if err != nil {
		panic(err)
	}

	if isAlive {
		fmt.Println("Session is alive.")
	} else {
		fmt.Println("Session is dead.")
		return
	}

	emailAddress, err := sess.GenerateEmailAddress()
	if err != nil {
		panic(err)
	}

	fmt.Println("Email address is " + emailAddress + ".")

	emails, err := sess.RetrieveMail(emailAddress)
	if err != nil {
		panic(err)
	}

	for _, email := range emails {
		fmt.Printf("From: %s, Subject: %s, Time: %s\n", email.From, email.Subject, email.Time)
	}
}
```
