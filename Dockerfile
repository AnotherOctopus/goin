FROM golang:1.9.7-stretch
WORKDIR /goinimage
ENV GIT_TERMINAL_PROMPT 1

ADD run/main.go /goinimage
COPY dump/ /goinimage/dump
ADD mongo.sh /goinimage
ADD networkfiles /goimage

RUN apt-get update
RUN apt-get -y install mongodb
RUN chmod 777 mongo.sh
RUN ./mongo.sh

ENV INSTALL_PATH $GOPATH/src/github.com/AnotherOctopus/goin
RUN go get gopkg.in/mgo.v2
RUN mkdir -p $INSTALL_PATH
COPY cnet $INSTALL_PATH/cnet
COPY wallet $INSTALL_PATH/wallet
COPY constants $INSTALL_PATH/constants

RUN echo $INSTALL_PATH

EXPOSE 1945
EXPOSE 1943
EXPOSE 1918
EXPOSE 1944

CMD ["go","run","main.go"]
