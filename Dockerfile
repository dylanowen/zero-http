FROM golang:latest AS build-env

ENV SRC_HOME "/go/src/"
ENV ZERO_HTTP_HOME "$SRC_HOME/github.com/dylanowen/zero-http"

# Copy our source over
ADD . $ZERO_HTTP_HOME

# compile zero-http
WORKDIR $ZERO_HTTP_HOME
# build go and bundle all the dependencies in the executable
RUN make publish-linux





FROM alpine:latest

ENV BUILD_HOME "/go/src/github.com/dylanowen/zero-http"
ENV ZERO_HOME "/zero"
ENV SRV_HOME "$ZERO_HOME/srv"
ENV CONFIG_HOME "$ZERO_HOME/.zero-http"

# setup our volumes for serving files
RUN mkdir -p $SRV_HOME
VOLUME ["$SRV_HOME"]

# setup our volumes for configuration
RUN mkdir -p $CONFIG_HOME
VOLUME ["$CONFIG_HOME"]

# Create a new group and user
RUN \
    addgroup -Sg 1000 zero &&  \
    adduser -SG zero -u 1000 -h $ZERO_HOME zero && \
    chown zero:zero $ZERO_HOME

WORKDIR $ZERO_HOME

# copy our build over
COPY --from=build-env $BUILD_HOME/zero-http /usr/local/bin

# Change to our zero-http user
USER zero

WORKDIR $SRV_HOME
EXPOSE 9000 9443