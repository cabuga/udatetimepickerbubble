package datetimepicker

const aiUsageDocument = `# udatetimepickerbubble

Bubble Tea datetime picker for editing correlated ISO week and calendar fields.
It supports full datetime selection and pure date selection.

## Import

` + "```go" + `
import datetimepicker "github.com/cabuga/udatetimepickerbubble"
` + "```" + `

## Picker Format

The component displays:

` + "```text" + `
CW 34.4  2026/08/20 22:21
` + "```" + `

In date mode it displays:

` + "```text" + `
CW 34.4  2026/08/20
` + "```" + `

Selectable fields are:

- calendar week: ` + "`34`" + `
- ISO weekday: ` + "`4`" + `
- year: ` + "`2026`" + `
- month: ` + "`08`" + `
- day: ` + "`20`" + `
- hour: ` + "`22`" + `
- minute: ` + "`21`" + `

Changing one field updates the others through Go time and ISO week rules.
For example, moving the day from 2026/08/20 to 2026/08/21 changes the display
from ` + "`CW 34.4`" + ` to ` + "`CW 34.5`" + `. Changing the calendar week from
` + "`CW 34.4`" + ` to ` + "`CW 36.4`" + ` changes the date to 2026/09/03.

## Controls

- left/right: select previous or next field
- up/down: increment or decrement selected field
- digits: replace selected field; commit happens when the field width is filled
- n: jump to the nearest next configured schedule
- t: set to the current time
- enter: confirm
- q or esc: cancel

## Basic Use

` + "```go" + `
package main

import (
	"fmt"
	"time"

	datetimepicker "github.com/carlos/udatetimepickerbubble"
)

func main() {
	picker := datetimepicker.New(datetimepicker.Config{
		InitialTime:  time.Now(),
		InitialField: datetimepicker.FieldCalendarWeek,
		ShowTutorial: true,
		Title:        "Select datetime",
	})

	result, err := picker.Run()
	if err != nil {
		panic(err)
	}
	if result.Canceled {
		return
	}

	fmt.Println(datetimepicker.Format(result.Time))
}
` + "```" + `

## Date-Only Use

` + "```go" + `
picker := datetimepicker.New(datetimepicker.Config{
	InitialTime:  time.Now(),
	InitialField: datetimepicker.FieldCalendarWeek,
	DateOnly:     true,
})

result, err := picker.Run()
if err != nil {
	panic(err)
}
fmt.Println(datetimepicker.FormatDate(result.Time))
` + "```" + `

## Schedule Jump

Pass one or more schedules and press ` + "`n`" + ` to jump to the nearest next
generated time. Cron expressions use five fields:

` + "```text" + `
minute hour day-of-month month day-of-week
` + "```" + `

For example, every Monday at 04:05:

` + "```go" + `
monday, err := datetimepicker.NewCronSchedule("5 4 * * 1")
if err != nil {
	panic(err)
}

picker := datetimepicker.New(datetimepicker.Config{
	InitialTime: time.Now(),
	Schedules: []datetimepicker.Schedule{
		monday,
	},
})
` + "```" + `

Multiple cron schedules are supported. The picker evaluates each schedule and
moves to the nearest generated time after the current picker value:

` + "```go" + `
monday, _ := datetimepicker.NewCronSchedule("5 4 * * 1")
thursday, _ := datetimepicker.NewCronSchedule("7 5 * * 4")

picker := datetimepicker.New(datetimepicker.Config{
	Schedules: []datetimepicker.Schedule{
		monday,
		thursday,
	},
})
` + "```" + `

If the picker points to Wednesday, pressing ` + "`n`" + ` selects Thursday at
05:07 because it is nearer than the following Monday at 04:05.

For a simple daily task at 21:00:

` + "```go" + `
picker := datetimepicker.New(datetimepicker.Config{
	InitialTime: time.Now(),
	Schedules: []datetimepicker.Schedule{
		datetimepicker.NewDailySchedule(21, 0),
	},
})
` + "```" + `

Use ` + "`ScheduleFunc`" + ` to adapt an external cron parser:

` + "```go" + `
picker := datetimepicker.New(datetimepicker.Config{
	Schedules: []datetimepicker.Schedule{
		datetimepicker.ScheduleFunc(func(after time.Time) time.Time {
			return cronSchedule.Next(after)
		}),
	},
})
` + "```" + `

Custom next-schedule keys can be configured:

` + "```go" + `
picker := datetimepicker.New(datetimepicker.Config{
	Schedules: []datetimepicker.Schedule{
		datetimepicker.NewDailySchedule(21, 0),
	},
	NextScheduleKeys: []string{"n", "N"},
})
` + "```" + `

The example CLI accepts repeated cron and daily schedules:

` + "```sh" + `
./bin/datetimepicker -cron '5 4 * * 1' -cron '7 5 * * 4'
./bin/datetimepicker -daily 21:00 -daily 23:30
` + "```" + `

## Use In An Existing Bubble Tea Program

` + "```go" + `
picker := datetimepicker.New(datetimepicker.Config{InitialTime: time.Now()})
model, cmd := picker.Update(msg)
picker = model.(datetimepicker.Model)
` + "```" + `

## go.work Use

From a parent workspace:

` + "```sh" + `
go work init
go work use ./path/to/your/app
go work use ./path/to/udatetimepickerbubble
` + "```" + `

Then import ` + "`github.com/carlos/udatetimepickerbubble`" + ` in your app. The
workspace will use the local checkout.

## Direct Git Use

` + "```sh" + `
go get github.com/carlos/udatetimepickerbubble
` + "```" + `
`

func AIUsageDocument() string {
	return aiUsageDocument
}
