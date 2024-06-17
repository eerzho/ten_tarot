package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eerzho/event_manager/internal/entity"
)

type AppleCalendar struct {
}

func NewAppleCalendar() *AppleCalendar {
	return &AppleCalendar{}
}

func (a *AppleCalendar) CreateFile(ctx context.Context, event *entity.Event) (string, error) {
	const op = "./internal/service/apple_calendar::CreateFile"

	filePath := filepath.Join(os.TempDir(), a.icsFilename(event))
	file, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	_, err = file.Write([]byte(a.icsContent(event)))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return filePath, nil
}

func (a *AppleCalendar) icsFilename(event *entity.Event) string {
	filename := fmt.Sprintf("%s, %s, %s, %s, %s, %s, %s, %s",
		event.Text,
		event.StartDate,
		event.EndDate,
		event.CTZ,
		event.Details,
		event.Location,
		event.CRM,
		event.Recur)

	hash := md5.Sum([]byte(filename))

	return fmt.Sprintf("%x", hash)
}

func (a *AppleCalendar) icsContent(event *entity.Event) string {
	var icsContent string

	if event.CTZ != "" {
		icsContent = fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VTIMEZONE
TZID:%s
BEGIN:STANDARD
DTSTART:20240101T020000
TZOFFSETFROM:+0000
TZOFFSETTO:+0000
TZNAME:GMT
END:STANDARD
END:VTIMEZONE
BEGIN:VEVENT
DTSTART;TZID=%s:%s
DTEND;TZID=%s:%s
SUMMARY:%s
DESCRIPTION:%s
STATUS:%s
`, event.CTZ, event.CTZ, event.StartDate, event.CTZ, event.EndDate, event.Text, event.Details, event.CRM)
	} else {
		icsContent = fmt.Sprintf(`BEGIN:VCALENDAR
VERSION:2.0
BEGIN:VEVENT
DTSTART:%s
DTEND:%s
SUMMARY:%s
DESCRIPTION:%s
STATUS:%s
`, event.StartDate, event.EndDate, event.Text, event.Details, event.CRM)
	}

	if event.Location != "" {
		icsContent += fmt.Sprintf("LOCATION:%s\n", event.Location)
	}

	if event.Recur != "" {
		if strings.Contains(event.Recur, "RRULE:") {
			icsContent += event.Recur + "\n"
		} else {
			icsContent += fmt.Sprintf("RRULE:%s\n", event.Recur)
		}
	}

	icsContent += "END:VEVENT\nEND:VCALENDAR"

	return icsContent
}
