FROM golang:1.19
WORKDIR /workspace/
COPY . .
RUN GOOS=linux go build -a -o printer -v ./main.go

FROM scratch
COPY --from=0 /workspace/printer /printer
COPY --from=0 /workspace/filament.tmpl /filament.tmpl

CMD ["/printer"]
