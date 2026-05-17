package utils

import (
    "net/smtp"
)

func SendEmail(to, subject, body string) error {
    from := "gnida3090@gmail.com"
    password := "ehyb wklc gyqv mizd"

    smtpHost := "smtp.gmail.com"
    smtpPort := "587"

    message := []byte(
        "Subject: " + subject + "\r\n" +
            "To: " + to + "\r\n" +
            "From: " + from + "\r\n\r\n" +
            body,
    )

    auth := smtp.PlainAuth("", from, password, smtpHost)

    err := smtp.SendMail(
        smtpHost+":"+smtpPort,
        auth,
        from,
        []string{to},
        message,
    )

    return err
}