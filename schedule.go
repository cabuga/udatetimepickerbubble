package datetimepicker

import "time"

type Schedule interface {
	Next(after time.Time) time.Time
}

type ScheduleFunc func(after time.Time) time.Time

func (f ScheduleFunc) Next(after time.Time) time.Time {
	if f == nil {
		return time.Time{}
	}

	return f(after)
}

func NewDailySchedule(hour int, minute int) Schedule {
	scheduleHour := clamp(hour, 0, 23)
	scheduleMinute := clamp(minute, 0, 59)

	return ScheduleFunc(func(after time.Time) time.Time {
		candidate := time.Date(
			after.Year(),
			after.Month(),
			after.Day(),
			scheduleHour,
			scheduleMinute,
			0,
			0,
			after.Location(),
		)
		if !candidate.After(after) {
			candidate = candidate.AddDate(0, 0, 1)
		}

		return candidate
	})
}
