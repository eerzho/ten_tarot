package service

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	"github.com/eerzho/ten_tarot/internal/entity"
)

type GoogleCalendar struct {
	url string
}

func NewGoogleCalendar(url string) *GoogleCalendar {
	return &GoogleCalendar{url: url}
}

func (g *GoogleCalendar) CreateUrl(ctx context.Context, event *entity.Event) string {
	const op = "./internal/service/google_calendar::CreateUrl"

	params := url.Values{}
	params.Add("action", "TEMPLATE")
	params.Add("text", event.Text)
	params.Add("dates", fmt.Sprintf("%s/%s", event.StartDate, event.EndDate))
	if event.CTZ != "" {
		params.Add("ctz", event.CTZ)
	}
	params.Add("details", event.Details)
	if event.Location != "" {
		params.Add("location", event.Location)
	}
	params.Add("crm", event.CRM)
	params.Add("trp", strconv.FormatBool(event.TRP))
	if event.Recur != "" {
		params.Add("recur", event.Recur)
	}

	return fmt.Sprintf("%s?%s", g.url, params.Encode())
}
