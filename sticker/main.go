package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"strings"

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

    if !sc.ValidateToken(os.Getenv("SLASH_TOKEN")) {
        log.Fatal("Token Error")
    }

    sticker := PickUpSticker(sc.Text)
    log.Printf("sticker: %s \n", sticker)

    image := GetStickerImage(sticker)
    log.Printf("image: %s \n", image)

    user := GetUserProfile(sc.UserID)
    log.Printf("user: %s, %s, %s \n", user.DisplayName, user.RealName, user.Image72)

    PostMessage(sc.ChannelName, user.RealName, user.Image72, image)
    log.Printf("Post Message \n")

	resp := Response{
		StatusCode:      204,
		IsBase64Encoded: false,
		Body:            "",
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

func PickUpSticker(str string) string {
    return strings.Replace(str, ":", "", -1)
}

func GetStickerImage(key string) string {
    api := slack.New(os.Getenv("LEGACY_TOKEN"))
    emojis, err := api.GetEmoji()
    if err != nil {
        log.Fatal("error: %s", err)
    }
    return emojis[key]
}

func GetUserProfile(key string) *slack.UserProfile {
    api := slack.New(os.Getenv("LEGACY_TOKEN"))
    profile, err := api.GetUserProfile(key, false)
    if err != nil {
        log.Fatal("Get Error: ", err)
    }
    return profile
}

func PostMessage(channel string, userName string, userIcon string, image string) {
    log.Printf("channel: %s, image: %s \n", channel, image)

    message := slack.WebhookMessage {
        Text: " ",
        Channel: "#" + channel,
        Username: userName,
        IconURL: userIcon,
        Attachments: []slack.Attachment {
            slack.Attachment {
                ImageURL: image,
            },
        },
    }

    err := slack.PostWebhook(os.Getenv("WEBHOOK_URL"), &message)
    if err != nil {
        log.Fatal("Get Error: ", err)
    }
}
