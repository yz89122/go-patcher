FROM golang:1

RUN go install 'github.com/yz89122/go-patcher@latest'

ADD example-app /src
ADD patches /src/patches

WORKDIR /src

RUN go build -overlay "$(go-patcher ./patches)" -o main ./main.go

CMD ./main
