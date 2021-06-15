FROM registry.access.redhat.com/ubi8/go-toolset:latest

ENV GOPATH=$APP_ROOT
ENV GOBIN=$APP_ROOT/bin

COPY . $GOPATH/src/mockda/
# This will be removed if we make the config as environment values
RUN go get gopkg.in/yaml.v2
RUN go install mockda

CMD mockda
