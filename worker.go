package main

import (
	"./JsonMessage"
	"./logger"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"net/smtp"
	"os"
	"runtime"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	//create logger
	str := "logfile" + JsonMessage.IntToString(int64(os.Getpid()))
	logger.CreateLogger(str)
	defer logger.DropLogger()

	//connect to RabbitMQ
	conn, err := amqp.Dial("amqp://admin:admin@stage.haru.io:5672/")
	logger.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	logger.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	msgs, err := ch.Consume(
		"email", // queue
		"",      // consumer   	consumer에 대한 식별자를 지정합니다. consumer tag는 로컬에 channel이므로, 두 클라이언트는 동일한 consumer tag를 사용할 수있다.
		false,   // autoAck    	false는 명시적 Ack를 해줘야 메시지가 삭제되고 true는 메시지를 빼면 바로 삭제
		false,   // exclusive	현재 connection에만 액세스 할 수 있으며, 연결이 종료 할 때 Queue가 삭제됩니다.
		false,   // noLocal    	필드가 설정되는 경우 서버는이를 published 연결로 메시지를 전송하지 않을 것입니다.
		false,   // noWait		설정하면, 서버는 Method에 응답하지 않습니다. 클라이언트는 응답 Method를 기다릴 것이다. 서버가 Method를 완료 할 수 없을 경우는 채널 또는 연결 예외를 발생시킬 것입니다.
		nil,     // arguments	일부 브로커를 사용하여 메시지의 TTL과 같은 추가 기능을 구현하기 위해 사용된다.
	)
	logger.FailOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			for d := range msgs {

				//Decoding arbitrary data
				var m JsonMessage.Email
				if err := json.Unmarshal([]byte(d.Body), &m); err != nil {
					logger.FailOnError(err, "Failed to json.Unmarshal")
					continue
				}

				body := "To: " + m.Address + "\r\nSubject: " +
					m.Title + "\r\n\r\n" + m.Body
				auth := smtp.PlainAuth("", "", "", "smtp.gmail.com")
				err := smtp.SendMail("smtp.gmail.com:587", auth, "",
					[]string{m.Address}, []byte(body))
				if err != nil {
					logger.FailOnError(err, "")
				}

				//RabbitMQ Message delete
				d.Ack(false)
			}

		}()
	}

	<-forever
	fmt.Printf("Worker is Dead.")

	os.Exit(0)
}
