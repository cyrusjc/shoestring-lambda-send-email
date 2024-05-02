package main

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

type Email struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

const (
	Sender    = "inquiry@shoestring.cafe"
	Recipient = "shoestringcafe@gmail.com"
	// The character encoding for the email.
	CharSet = "UTF-8"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// field to fill struct is in body
	body := request.Body
	var email Email

	log.Println("Messaged received:", body)
	err := json.Unmarshal([]byte(body), &email)
	if err != nil {
		log.Println("Error unmarshaling: ", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Unmarshaling body" + err.Error()),
			StatusCode: 400,
		}, nil
	}

	email.Message = strings.TrimSuffix(email.Message, "\n")

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	emailClient := ses.New(sess)

	data := fmt.Sprintf("{ \"name\":\""+"%s"+"\", \"phone\": \""+"%s"+"\", \"email\": \""+"%s"+"\", \"message\": \""+"%s"+"\"}", email.Name, email.Phone, email.Email, email.Message)
	input := &ses.SendTemplatedEmailInput{
		Source:           aws.String(Sender),
		Template:         aws.String("shoestringEmailTemplate"),
		ReplyToAddresses: []*string{aws.String(email.Email)},
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(Recipient)},
		},
		TemplateData: aws.String(data),
	}
	_, err = emailClient.SendTemplatedEmail(input)

	if err != nil {
		log.Println("Error sending email: ", err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error sending email" + err.Error()),
			StatusCode: 400,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		Body:       "Success",
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(Handler)
}
