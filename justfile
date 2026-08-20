build:
    #here put here the commit

build-example:
    mkdir -p bin
    go build -buildvcs=false -o bin/datetimepicker ./cmd/datetimepicker
    @echo "Example built: ./bin/datetimepicker"
    @echo "Run picker: ./bin/datetimepicker -time '2026-08-20 22:21' -field cw"
    @echo "AI doc: ./bin/datetimepicker -aigeneration"

revive:
    $(go env GOPATH)/bin/revive -config "$PWD/revive.toml" ./...
