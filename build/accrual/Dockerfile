FROM alpine:latest

RUN apk --no-cache add libc6-compat

COPY ./cmd/accrual/accrual_linux_amd64 /bin

RUN chmod +x /bin/accrual_linux_amd64

CMD ["/bin/accrual_linux_amd64"]
