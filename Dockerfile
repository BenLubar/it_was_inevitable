FROM golang:1.13 as builder

COPY go.mod go.sum /src/it_was_inevitable/

WORKDIR /src/it_was_inevitable

RUN go mod download

COPY *.go /src/it_was_inevitable/

RUN CGO_ENABLED=0 go build -o /it_was_inevitable

FROM benlubar/dwarffortress:df-ai-0.44.12-r2

COPY --from=builder /it_was_inevitable /usr/local/bin/it_was_inevitable

RUN sed -i /df_linux/dfhack -e "s/ setarch / exec setarch /"

ENTRYPOINT ["it_was_inevitable"]
