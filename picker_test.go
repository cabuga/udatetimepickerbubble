package datetimepicker

import (
	"testing"
	"time"

	tea "charm.land/bubbletea/v2"
)

func TestFormatInitialExample(t *testing.T) {
	value := time.Date(2026, time.August, 20, 22, 21, 44, 0, time.UTC)

	got := Format(value)
	want := "CW 34.4  2026/08/20 22:21"

	if got != want {
		t.Fatalf("Format() = %q, want %q", got, want)
	}
}

func TestDayAdjustUpdatesISOWeekday(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		InitialField: FieldDay,
	})

	updated := updateWithKey(t, model, "up")

	got := Format(updated.CurrentTime())
	want := "CW 34.5  2026/08/21 22:21"
	if got != want {
		t.Fatalf("after day up = %q, want %q", got, want)
	}
}

func TestCalendarWeekAdjustUpdatesDate(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 22, 21, 0, 0, time.UTC),
		InitialField: FieldCalendarWeek,
	})

	updated := updateWithKey(t, model, "3")
	updated = updateWithKey(t, updated, "6")

	got := Format(updated.CurrentTime())
	want := "CW 36.4  2026/09/03 22:21"
	if got != want {
		t.Fatalf("after week 36 = %q, want %q", got, want)
	}
}

func TestMonthChangeClampsDay(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.March, 31, 10, 30, 0, 0, time.UTC),
		InitialField: FieldMonth,
	})

	updated := updateWithKey(t, model, "down")

	got := Format(updated.CurrentTime())
	want := "CW 09.6  2026/02/28 10:30"
	if got != want {
		t.Fatalf("after month down = %q, want %q", got, want)
	}
}

func TestDayDownMovesToPreviousMonth(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.March, 1, 10, 30, 0, 0, time.UTC),
		InitialField: FieldDay,
	})

	updated := updateWithKey(t, model, "down")

	got := Format(updated.CurrentTime())
	want := "CW 09.6  2026/02/28 10:30"
	if got != want {
		t.Fatalf("after day down = %q, want %q", got, want)
	}
}

func TestFormatDateOmitsTime(t *testing.T) {
	value := time.Date(2026, time.August, 20, 23, 21, 0, 0, time.UTC)

	got := FormatDate(value)
	want := "CW 34.4  2026/08/20"

	if got != want {
		t.Fatalf("FormatDate() = %q, want %q", got, want)
	}
}

func TestDateModeSkipsTimeFields(t *testing.T) {
	model := New(Config{
		InitialTime:  time.Date(2026, time.August, 20, 23, 21, 0, 0, time.UTC),
		InitialField: FieldDay,
		DateOnly:     true,
	})

	updated := updateWithKey(t, model, "right")

	if updated.SelectedField() != FieldCalendarWeek {
		t.Fatalf(
			"selected field = %v, want %v",
			updated.SelectedField(),
			FieldCalendarWeek,
		)
	}

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 00:00"
	if got != want {
		t.Fatalf("date mode current time = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesDailySchedule(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(2026, time.August, 20, 20, 30, 0, 0, time.UTC),
		Schedules:   []Schedule{NewDailySchedule(21, 0)},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 21:00"
	if got != want {
		t.Fatalf("next schedule = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesNearestSchedule(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(2026, time.August, 20, 20, 30, 0, 0, time.UTC),
		Schedules: []Schedule{
			NewDailySchedule(23, 0),
			NewDailySchedule(21, 0),
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 21:00"
	if got != want {
		t.Fatalf("nearest schedule = %q, want %q", got, want)
	}
}

func TestNextScheduleCanUseCustomKey(t *testing.T) {
	model := New(Config{
		InitialTime: time.Date(
			2026,
			time.August,
			20,
			21,
			0,
			0,
			0,
			time.UTC,
		),
		Schedules: []Schedule{NewDailySchedule(21, 0)},
		NextScheduleKeys: []string{
			"x",
		},
	})

	updated := updateWithKey(t, model, "x")

	got := Format(updated.CurrentTime())
	want := "CW 34.5  2026/08/21 21:00"
	if got != want {
		t.Fatalf("custom key next schedule = %q, want %q", got, want)
	}
}

func TestCronScheduleNextMonday(t *testing.T) {
	schedule, err := NewCronSchedule("5 4 * * 1")
	if err != nil {
		t.Fatalf("NewCronSchedule() error = %v", err)
	}

	model := New(Config{
		InitialTime: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		Schedules: []Schedule{
			schedule,
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 35.1  2026/08/24 04:05"
	if got != want {
		t.Fatalf("cron next monday = %q, want %q", got, want)
	}
}

func TestNextScheduleUsesNearestCronSchedule(t *testing.T) {
	monday, mondayErr := NewCronSchedule("5 4 * * 1")
	if mondayErr != nil {
		t.Fatalf("NewCronSchedule(monday) error = %v", mondayErr)
	}
	thursday, thursdayErr := NewCronSchedule("7 5 * * 4")
	if thursdayErr != nil {
		t.Fatalf("NewCronSchedule(thursday) error = %v", thursdayErr)
	}

	model := New(Config{
		InitialTime: time.Date(2026, time.August, 19, 10, 0, 0, 0, time.UTC),
		Schedules: []Schedule{
			monday,
			thursday,
		},
	})

	updated := updateWithKey(t, model, "n")

	got := Format(updated.CurrentTime())
	want := "CW 34.4  2026/08/20 05:07"
	if got != want {
		t.Fatalf("nearest cron = %q, want %q", got, want)
	}
}

func updateWithKey(t *testing.T, model Model, key string) Model {
	t.Helper()

	var updated tea.Model
	switch key {
	case "left":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyLeft}))
	case "right":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyRight}))
	case "up":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyUp}))
	case "down":
		updated, _ = model.Update(tea.KeyPressMsg(tea.Key{Code: tea.KeyDown}))
	default:
		keyMsg := tea.KeyPressMsg(tea.Key{Code: []rune(key)[0], Text: key})
		updated, _ = model.Update(keyMsg)
	}

	pickerModel, ok := updated.(Model)
	if !ok {
		t.Fatalf("updated model type = %T, want datetimepicker.Model", updated)
	}

	return pickerModel
}
