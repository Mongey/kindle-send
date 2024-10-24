package mail

import (
	"log"
	"os"

	"github.com/nikhil1raghav/kindle-send/util"

	config "github.com/nikhil1raghav/kindle-send/config"

	mail "github.com/wneessen/go-mail"
)

func Send(files []string, timeout int) {
	cfg := config.GetInstance()
	message := mail.NewMsg()
	if err := message.From(cfg.Sender); err != nil {
		util.Red.Printf("failed to set From address: %s", err)
		return
	}
	if err := message.To(cfg.Receiver); err != nil {
		util.Red.Printf("failed to set To address: %s", err)
		return
	}
	message.Subject("Book!")
	message.SetBodyString(mail.TypeTextPlain, "See attached files.")

	attachedFiles := make([]string, 0)
	for _, file := range files {
		_, err := os.Stat(file)
		if err != nil {
			util.Red.Printf("Couldn't find the file %s : %s \n", file, err)
			continue
		} else {
			message.AttachFile(file)
			attachedFiles = append(attachedFiles, file)
		}
	}
	if len(attachedFiles) == 0 {
		util.Cyan.Println("No files to send")
		return
	}

	util.CyanBold.Println("Sending mail")
	util.Cyan.Println("Following files will be sent :")
	for i, file := range attachedFiles {
		util.Cyan.Printf("%d. %s\n", i+1, file)
	}

	client, err := mail.NewClient(
		cfg.Server,

		mail.WithSMTPAuth(mail.SMTPAuthPlain),
		mail.WithTLSPortPolicy(mail.TLSMandatory),
		mail.WithUsername(cfg.Sender),
		mail.WithPassword(cfg.Password),
	)

	if err != nil {
		util.Red.Printf("failed to construct client: %s\n", err)
	}

	if err := client.DialAndSend(message); err != nil {
		util.Red.Printf("failed to send mail: %s\n", err)
	} else {
		util.GreenBold.Printf("Mailed %d files to %s", len(attachedFiles), cfg.Receiver)
	}

}
