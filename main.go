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

func main() {
	lambda.Start(Handler)

}

const (
	// Replace sender@example.com with your "From" address.
	// This address must be verified with Amazon SES.
	Sender = "shoestringcafe@gmail.com"

	// Replace recipient@example.com with a "To" address. If your account
	// is still in the sandbox, this address must be verified.
	Recipient = "shoestringcafe@gmail.com"

	// The character encoding for the email.
	CharSet = "UTF-8"
)

type Email struct {
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Email   string `json:"email"`
	Message string `json:"message"`
}

func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

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

	emailTemplate := "Name:" + email.Name + "\n" + "Email: " + email.Email + "\n" + "Phone: " + email.Phone + "\n" + "Body:" + email.Message

	input := &ses.SendEmailInput{
		Destination: &ses.Destination{
			CcAddresses: []*string{
				// aws.String(email.Email),
			},
			ToAddresses: []*string{
				aws.String(Recipient),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailTemplate),
				},
				Text: &ses.Content{
					Charset: aws.String(CharSet),
					Data:    aws.String(emailTemplate),
				},
			},
			Subject: &ses.Content{
				Charset: aws.String(CharSet),
				Data:    aws.String("[Shoestring Cafe]" + email.Name + " " + email.Email),
			},
		},
		Source: aws.String(Sender),
	}

	_, err = emailClient.SendEmail(input)

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
