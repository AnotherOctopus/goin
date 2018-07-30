FROM golang:1.9.7-stretch
WORKDIR /goinimage
ADD run/main.go /goinimage
ENV GIT_TERMINAL_PROMPT 1
ENV INSTALL_PATH $GOPATH/src/github.com/AnotherOctopus/goin

RUN mkdir -p $INSTALL_PATH
RUN go get gopkg.in/mgo.v2

COPY cnet $INSTALL_PATH/cnet
COPY wallet $INSTALL_PATH/wallet
COPY constants $INSTALL_PATH/constants

RUN echo $INSTALL_PATH

EXPOSE 1945
EXPOSE 1943
EXPOSE 1918
EXPOSE 1944
CMD ["go","run","main.go"]
