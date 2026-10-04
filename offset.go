package datetimepicker

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const offsetPrefix = "offset:"
const shortOffsetPrefix = "o:"
const cronPrefix = "c:"

var offsetPartPattern = regexp.MustCompile(`([0-9]+(?:\.[0-9]+)?)(ns|us|µs|μs|ms|s|m|h|d|w)`)

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
	duration, err := parseOffsetDuration(durationText)
	if err != nil {
		return nil, fmt.Errorf("parse offset duration %q: %w", durationText, err)
	}
	if duration <= 0 {
		return nil, fmt.Errorf("offset duration must be positive")
	}
	return offsetSchedule{duration: duration}, nil
}

func parseOffsetDuration(raw string) (time.Duration, error) {
	text := strings.TrimPrefix(raw, "+")
	parts := offsetPartPattern.FindAllStringSubmatchIndex(text, -1)
	if len(parts) == 0 {
		return 0, fmt.Errorf("invalid duration")
	}
	duration := time.Duration(0)
	position := 0
	for _, part := range parts {
		if part[0] != position {
			return 0, fmt.Errorf("invalid duration")
		}
		valueText := text[part[2]:part[3]]
		unit := text[part[4]:part[5]]
		var component time.Duration
		if unit == "d" || unit == "w" {
			value, err := strconv.ParseFloat(valueText, 64)
			if err != nil {
				return 0, err
			}
			factor := float64(24 * time.Hour)
			if unit == "w" {
				factor *= 7
			}
			nanoseconds := value * factor
			if math.IsInf(nanoseconds, 0) || nanoseconds > math.MaxInt64 {
				return 0, fmt.Errorf("duration overflows time.Duration")
			}
			component = time.Duration(nanoseconds)
		} else {
			var err error
			component, err = time.ParseDuration(text[part[0]:part[1]])
			if err != nil {
				return 0, err
			}
		}
		if component > 0 && duration > time.Duration(math.MaxInt64)-component {
			return 0, fmt.Errorf("duration overflows time.Duration")
		}
		duration += component
		position = part[1]
	}
	if position != len(text) {
		return 0, fmt.Errorf("invalid duration")
	}
	return duration, nil
}

// NewSchedule parses a cron or offset expression. Cron expressions can use
// the c: prefix or be bare; offset expressions can use o: or offset:.
func NewSchedule(expression string) (Schedule, error) {
	expression = strings.TrimSpace(expression)
	if strings.HasPrefix(expression, cronPrefix) {
		return NewCronSchedule(strings.TrimPrefix(expression, cronPrefix))
	}
	if strings.HasPrefix(expression, shortOffsetPrefix) {
		return NewOffsetSchedule(offsetPrefix + strings.TrimPrefix(expression, shortOffsetPrefix))
	}
	if strings.HasPrefix(expression, offsetPrefix) {
		return NewOffsetSchedule(expression)
	}
	return NewCronSchedule(expression)
}
