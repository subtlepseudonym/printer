FROM golang:1.19
WORKDIR /workspace/
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -o printer -v ./main.go

FROM scratch
COPY --from=0 /workspace/printer /printer
COPY --from=0 /workspace/filament.tmpl /filament.tmpl
COPY --from=tarampampam/curl:latest /bin/curl /curl

CMD ["/printer"]
