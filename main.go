package main

import (
	"fmt"
	"net/http"
	"os"

	goteamsnotify "github.com/atc0005/go-teams-notify/v2"
	"github.com/labstack/echo/v4"
	"github.com/philips-software/go-hsdp-api/console/metrics/alerts"
	"github.com/spf13/viper"
)

func main() {
	viper.SetEnvPrefix("events")
	viper.SetDefault("token", "")
	viper.SetDefault("webhook_url", "")
	viper.AutomaticEnv()

	token := viper.GetString("token")
	if token == "" {
		fmt.Printf("empty token\n")
		os.Exit(1)
	}
	webhookURL := viper.GetString("webhook_url")
	if webhookURL == "" {
		fmt.Printf("empty webhook_url\n")
		os.Exit(2)
	}

	e := echo.New()
	e.POST("/webhook/:token", webhookProcessor(token, webhookURL))

	_ = e.Start(":8080")
}

type responseMessage struct {
	Message string `json:"message"`
}

func webhookProcessor(token, webhookURL string) echo.HandlerFunc {
	client := goteamsnotify.NewClient()

	return func(ctx echo.Context) error {
		var payload alerts.Payload

		if ctx.Param("token") != token {
			return ctx.JSON(http.StatusUnauthorized, responseMessage{"unauthorized"})
		}
		if err := ctx.Bind(&payload); err != nil {
			return ctx.JSON(http.StatusBadRequest, responseMessage{err.Error()})
		}
		err := convertAndSendToTeams(client, webhookURL, payload)
		if err != nil {
			return ctx.JSON(http.StatusInternalServerError, responseMessage{err.Error()})
		}
		return ctx.JSON(http.StatusOK, responseMessage{"delivered"})
	}
}

func convertAndSendToTeams(client goteamsnotify.API, webhookURL string, payload alerts.Payload) error {
	msgCard := goteamsnotify.NewMessageCard()
	switch payload.Status {
	case "firing":
		msgCard.ThemeColor = "d71900"
	case "resolved":
		msgCard.ThemeColor = "0076D7"
	default:
		msgCard.ThemeColor = "d79700"
	}
	msgCard.Title = payload.CommonAnnotations.Description
	msgCard.Summary = payload.CommonAnnotations.Summary
	labelSection := goteamsnotify.NewMessageCardSection()
	labelSection.ActivityTitle = "Details"
	_ = labelSection.AddFact(goteamsnotify.MessageCardSectionFact{
		Name:  "application",
		Value: payload.CommonLabels.Application,
	},
		goteamsnotify.MessageCardSectionFact{
			Name:  "applicationId",
			Value: payload.CommonLabels.ApplicationID,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "brokerId",
			Value: payload.CommonLabels.BrokerID,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "region",
			Value: payload.CommonLabels.Deployment,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "organization",
			Value: payload.CommonLabels.Organization,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "space",
			Value: payload.CommonLabels.Space,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "job",
			Value: payload.CommonLabels.Job,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "severity",
			Value: payload.CommonLabels.Severity,
		},
		goteamsnotify.MessageCardSectionFact{
			Name:  "instance",
			Value: payload.CommonLabels.Instance,
		},
	)
	_ = msgCard.AddSection(labelSection)

	return client.Send(webhookURL, msgCard)
}
