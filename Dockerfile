# Go Test DB Worker
#
# VERSION 0.0.1

FROM ubuntu_14.04_w_git:latest

MAINTAINER Unai Garcia <unai@gamewheel.com>
LABEL version="0.0.1" description="Go DB Worker App which connects to DB and communicates thru rabbitmq"

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		curl \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.5.3
ENV GOLANG_DOWNLOAD_URL https://golang.org/dl/go$GOLANG_VERSION.linux-amd64.tar.gz
ENV GOLANG_DOWNLOAD_SHA256 43afe0c5017e502630b1aea4d44b8a7f059bf60d7f29dfd58db454d4e4e0ae53

RUN curl -fsSLk "$GOLANG_DOWNLOAD_URL" -o golang.tar.gz \
	&& echo "$GOLANG_DOWNLOAD_SHA256  golang.tar.gz" | sha256sum -c - \
	&& tar -C /usr/local -xzf golang.tar.gz \
	&& rm golang.tar.gz

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"

ENV GO_TEST_DB_WORKER_PATH=/go/src/github.com/ugarcia/go_test_db_worker

RUN mkdir -p $GO_TEST_DB_WORKER_PATH
WORKDIR $GO_TEST_DB_WORKER_PATH

EXPOSE 8601
CMD ["go", "get"]
CMD ["go", "run", "main.go"]
