package startup

import (
	"github.com/getsentry/sentry-go"
	"log"
	"madaurus/dev/material/app/interfaces"
	"net/http"
)

func InitSentry(SentrySetting *interfaces.Sentry) {
	sentry.Init(sentry.ClientOptions{
		Dsn:           SentrySetting.DNS,
		EnableTracing: SentrySetting.EnableTracing,
		BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
			if hint.Context != nil {
				if _, ok := hint.Context.Value(sentry.RequestContextKey).(*http.Request); ok {
					log.Printf("Request data is attached to the event")
				}
			}
			return event
		},
		Debug: SentrySetting.DEBUG,
	})
}
