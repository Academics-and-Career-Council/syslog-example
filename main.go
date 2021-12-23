package main

import (
	"fmt"

	"github.com/streadway/amqp"
	"encoding/json"
	"github.com/go-co-op/gocron"
	"time"
)

func main() {
	// ****** USE THIS FUNCTION TO WRITE LOGS TO RABBITMQ ******
	

	fmt.Println("Go RabbitMQ Tutorial")
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	// This function connects to rabbitmq. The URL is similar to that of mongoDB url,
	// but difference is that here it is connected to a port in the local system.
	// To break down the URL, "amqp" is standard notation, guest:guest is username and password
	// which is sort of standard, and 5672 is the port. You can check it out by searching
	// this port in the browser and see if the port is running 
	if err != nil {
		fmt.Println(err)
		panic(1)
	}
	defer conn.Close()

	fmt.Println("Successfully Connected to our RabbitMQ Instance")

	ch, err := conn.Channel()// This creates a channel to send data.
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"TestQueue",
		false,
		false,
		false,
		false,
		nil,
	)// These are standard values, not to be changed(except the queue name,
		// which has to be same as that of the reciever).
	// Basically declarations to the queue

	if err != nil  {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println(q)

	type Syslog struct {
		ServiceName string		`bson:"service_name"`
		StatusCode	int			`bson:"status_code"`
		Severity	string		`bson:"severity"`
		MsgName		string		`bson:"msg_name"`
		Msg			string		`bson:"msg"`
		InvokedBy	string		`bson:"invoked_by"`
		Result		string		`bson:"result"`
		Batch      	int    		`bson:"batch"`
		Timestamp	time.Time	`bson:"timestamp"`
		CreatedAt 	time.Time   `bson:"createdAt,omitempty"`
	}//Object to be sent to syslog, should follow the same pattern
	log := Syslog{
		"xenon",
		409,
		"error", 
		"xenon error", 
		"check xenon", 
		"student@iitk.ac.in", 
		"failure", 
		-2, 
		time.Now(), 
		time.Now(),
	}

	data, _ := json.Marshal(log)//convert the data to []byte format, 
			//so it can be easily converted to json and read by syslog
	BodyJson := string(data)//converts the []byte to string, the only way 
			//to send via RabbitMQ

	err = ch.Publish(
		"",
		"TestQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(BodyJson),
		},
	)//publishes a result to the queue
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	fmt.Println("successfully published message to queue")
	//****** COVERS BASIC SENDING TO RABBITMQ WHICH THEN SENDS IT 
	//         TO MONGODB USING THE LOG APPLICATION ******
	//**** Sends multiple data periodically ****
	s := gocron.NewScheduler(time.UTC)
	s.Every(5).Seconds().Do(func(){
		log.CreatedAt = time.Now()
		log.Timestamp = time.Now()
		data, _ := json.Marshal(log)
		BodyJson := string(data)
		err = ch.Publish(
			"",
			"TestQueue",
			false,
			false,
			amqp.Publishing{
				ContentType: "text/plain",
				Body: []byte(BodyJson),
			},
		)
		if err != nil {
			fmt.Println(err)
			panic(err)
		}
		fmt.Println("successfully published message to queue")
	})
	s.StartAsync()
	s.StartBlocking()
	//**** end ****
}
