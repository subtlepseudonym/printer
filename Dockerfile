FROM golang:1.19
WORKDIR /workspace/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o printer -v ./main.go

FROM scratch
COPY --from=0 /workspace/printer /printer
COPY --from=0 /workspace/filament.tmpl /filament.tmpl

CMD ["/printer"]
