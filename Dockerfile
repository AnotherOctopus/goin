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
RUN go get gopkg.in/mgo.v2
RUN go get github.com/gorilla/mux
RUN go get github.com/gorilla/handlers
RUN mkdir -p $INSTALL_PATH
COPY cnet $INSTALL_PATH/cnet
COPY wallet $INSTALL_PATH/wallet
COPY constants $INSTALL_PATH/constants
COPY network $INSTALL_PATH/network

RUN echo $INSTALL_PATH

EXPOSE 1945
EXPOSE 1943
EXPOSE 1918
EXPOSE 1944
EXPOSE 80

CMD ["./mongo.sh"]
#CMD ["hostname"]
