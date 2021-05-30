# --- builder ---
# golang:1.16.4-alpine3.13
FROM golang@sha256:9dd1788d4bd0df3006d79a88bda67cb8357ab49028eebbcb1ae64f2ec07be627 AS build
COPY . .
RUN unset GOPATH && \
    CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o /go/bin/tfswitch

# --- realease ---
# alpine3.13
FROM alpine@sha256:def822f9851ca422481ec6fee59a9966f12b351c62ccb9aca841526ffaa9f748 AS release
COPY --from=build /go/bin/tfswitch /usr/bin/tfswitch
# tfswitch -u: Testing functionality by installing Terraform
# Since the Terraform binary is less than 100b, it has a minimal effect on final image size
# tfstich -v: this is not a good test and it does not actually check functionality 
RUN tfswitch -u
ENTRYPOINT [ "tfswitch" ]
