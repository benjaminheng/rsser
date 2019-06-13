FROM golang:1.12-alpine

RUN apk add git

COPY . /code
VOLUME ["/code"]
RUN cd /code/cmd/server && go install .
EXPOSE 8000
CMD ["server"]
