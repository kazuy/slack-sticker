package main

import (
	"context"
	"log"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joncalhoun/qson"
	"github.com/slack-go/slack"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (Response, error) {
    sc := QueryParameterToSlashCommand(request.Body)

    msg := Message(sc)

	resp := Response{
		StatusCode:      200,
		IsBase64Encoded: false,
		Body:            msg,
		Headers: map[string]string{
			"Content-Type":           "application/json",
			"X-MyCompany-Func-Reply": "sticker-handler",
		},
	}

	return resp, nil
}

func main() {
	lambda.Start(Handler)
}

func QueryParameterToSlashCommand(q string) *slack.SlashCommand {
    log.Printf("Body: %s\n", q)

    b, err := qson.ToJSON(q)
    if err != nil {
        log.Fatal("Convert failed: ", err)
    }

    s := new(slack.SlashCommand)
    err = json.Unmarshal(b, s)
    if err != nil {
        log.Fatal("Convert failed: ", err)
    }

    return s
}

func Message(sc *slack.SlashCommand) string {
    params := &slack.Msg {
        Text: sc.Text,
    }

    b, err := json.Marshal(params)
    if err != nil {
        log.Fatal("Convert failed: ", err)
    }

    return string(b)
}
