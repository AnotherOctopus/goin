FROM golang:1.9.7-stretch
WORKDIR /goinimage
ADD run/main.go /goinimage
ENV GIT_TERMINAL_PROMPT 1
RUN go get github.com/AnotherOctopus/goin/wallet/
RUN go get github.com/AnotherOctopus/goin/cnet/
RUN go get github.com/AnotherOctopus/goin/constants/
EXPOSE 1945
EXPOSE 1943
EXPOSE 1918
EXPOSE 1944
CMD ["go","run","main.go"]
