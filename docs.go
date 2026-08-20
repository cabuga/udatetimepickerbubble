package datetimepicker

const aiUsageDocument = `# udatetimepickerbubble

Bubble Tea datetime picker for editing correlated ISO week and calendar fields.

## Import

` + "```go" + `
import datetimepicker "github.com/carlos/udatetimepickerbubble"
` + "```" + `

## Picker Format

The component displays:

` + "```text" + `
CW 34.4  2026/08/20 22:21
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
