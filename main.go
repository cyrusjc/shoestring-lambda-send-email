package main

import (
	"encoding/json"
	"fmt"

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
	Sender    = "shoestringcafe@gmail.com"
	Recipient = "shoestringcafe@gmail.com"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	// field to fill struct is in body
	body := request.Body
	var email Email

	err := json.Unmarshal([]byte(body), &email)
	if err != nil {
		fmt.Println(err)
		return events.APIGatewayProxyResponse{
			Body:       fmt.Sprintf("Error Unmarshaling body" + err.Error()),
			StatusCode: 400,
		}, nil
	}

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("us-west-2")},
	)
	emailClient := ses.New(sess)

	data := "{ \"name\":\"" + email.Name + "\", \"phone\": \"" + email.Phone + "\", \"email\": \"" + email.Email + "\", \"message\": \"" + email.Message + "\"}"
	input := &ses.SendTemplatedEmailInput{
		Source:   aws.String(Sender),
		Template: aws.String("shoestringEmailTemplate"),
		Destination: &ses.Destination{
			ToAddresses: []*string{aws.String(Recipient)},
		},
		TemplateData: &data,
	}
	_, err = emailClient.SendTemplatedEmail(input)

	if err != nil {
		fmt.Println(err)
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
