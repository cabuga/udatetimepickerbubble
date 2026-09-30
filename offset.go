package datetimepicker

import (
	"fmt"
	"strings"
	"time"
)

const offsetPrefix = "offset:"

type offsetSchedule struct {
	duration time.Duration
}

func (s offsetSchedule) Next(after time.Time) time.Time {
	return after.Add(s.duration)
}

// NewOffsetSchedule parses an offset expression such as "offset:+72h".
func NewOffsetSchedule(expression string) (Schedule, error) {
	if !strings.HasPrefix(expression, offsetPrefix) {
		return nil, fmt.Errorf("offset expression %q must start with %q", expression, offsetPrefix)
	}
	durationText := strings.TrimPrefix(expression, offsetPrefix)
	if durationText == "" {
		return nil, fmt.Errorf("offset expression %q has an empty duration", expression)
	}
	duration, err := time.ParseDuration(durationText)
	if err != nil {
		return nil, fmt.Errorf("parse offset duration %q: %w", durationText, err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("offset duration must be positive")
	}
	return offsetSchedule{duration: duration}, nil
}

// NewSchedule parses an offset expression, or a cron expression when no
// offset: prefix is present.
func NewSchedule(expression string) (Schedule, error) {
	if strings.HasPrefix(expression, offsetPrefix) {
		return NewOffsetSchedule(expression)
	}
	return NewCronSchedule(expression)
}
