FROM golang:1.11.2 as builder

COPY *.go /go/src/github.com/BenLubar/it_was_inevitable/

WORKDIR /go/src/github.com/BenLubar/it_was_inevitable

RUN go get -d

ARG tag=
RUN CGO_ENABLED=0 go build -a -tags "$tag" -o /it_was_inevitable

FROM benlubar/dwarffortress:df-ai-0.44.12-r1-update1

COPY --from=builder /it_was_inevitable /usr/local/bin/it_was_inevitable

RUN sed -i /df_linux/dfhack -e "s/ setarch / exec setarch /"

CMD ["it_was_inevitable"]
