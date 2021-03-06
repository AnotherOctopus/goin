FROM golang:1.9.7-stretch
WORKDIR /goinimage
ENV GIT_TERMINAL_PROMPT 1

ADD run/ /goinimage/run
ADD mongo.sh /goinimage
ADD dump/ /goinimage/dump
ADD networkfiles /goinimage/networkfiles

RUN apt-get update
RUN apt-get -y install mongodb
RUN chmod 777 mongo.sh

ENV INSTALL_PATH $GOPATH/src/github.com/AnotherOctopus/goin
RUN go get  github.com/globalsign/mgo
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN mkdir -p $INSTALL_PATH
COPY cnet $INSTALL_PATH/cnet
COPY wallet $INSTALL_PATH/wallet
COPY constants $INSTALL_PATH/constants
COPY network $INSTALL_PATH/network

RUN echo $INSTALL_PATH


CMD ["./mongo.sh"]
