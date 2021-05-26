FROM alpine:3.11 AS build
RUN apk upgrade -U -a && \
          apk upgrade && \
          apk add --update go gcc g++ git ca-certificates curl make 
WORKDIR /app
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build && mv terraform-switcher /usr/local/bin/tfswitch
RUN echo $PATH
ENTRYPOINT [ "tfswitch" ]

# IF tfswitch command does not work, the buid will exit with non-zero and CI will stop
FROM build as test
RUN tfswitch -v
