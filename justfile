build:
    #here put here the commit

build-example:
    mkdir -p bin
    GOCACHE=/tmp/go-build-cache go build -buildvcs=false -o bin/datetimepicker ./cmd/datetimepicker
    @echo "Example built: ./bin/datetimepicker"
    @echo "Datetime: ./bin/datetimepicker -time '2026-08-20 22:21' -field cw"
    @echo "Date: ./bin/datetimepicker -mode date -time '2026-08-20'"
    @echo "Schedules: ./bin/datetimepicker -daily 21:00 -daily 23:30"
    @echo "Cron: ./bin/datetimepicker -cron '5 4 * * 1' -cron '7 5 * * 4'"
    @echo "AI doc: ./bin/datetimepicker -aigeneration"

example: build-example

revive:
    $(go env GOPATH)/bin/revive -config "$PWD/revive.toml" ./...
