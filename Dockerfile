FROM golang:1.23-alpine AS build

WORKDIR /build

# copy source code
COPY ./cmd/ ./cmd/
COPY ./go.mod ./
COPY ./go.sum ./

# build the executable
RUN go build -o webhooker ./cmd

# done with the build stage, onto the deploy stage
FROM alpine:latest

WORKDIR /app

# copy the built executable
COPY --from=build /build/webhooker /app/webhooker

# set the default config path
ENV WEBHOOKER_CONFIG_PATH="/app/webhooker.yml"
EXPOSE 80 443

ENTRYPOINT [ "/app/webhooker" ]