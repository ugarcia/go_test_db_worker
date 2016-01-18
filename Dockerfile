# Go Test DB Worker
#
# VERSION 0.0.1

FROM golang:latest

MAINTAINER Unai Garcia <unai@gamewheel.com>
LABEL version="0.0.1" description="Go DB Worker App which connects to DB and communicates thru rabbitmq"

# Copy git ssh key
ENV HOME /root
RUN mkdir -p $HOME/.ssh
ADD id_rsa $HOME/.ssh/id_rsa
RUN chmod 700 $HOME/.ssh/id_rsa && \
          echo "Host github.com\n\tStrictHostKeyChecking no\n" >> $HOME/.ssh/config && \
          git config --global url.ssh://git@github.com/.insteadOf https://github.com/

# Create app folder, copy source and set as working dir
RUN mkdir -p /go/src/app
COPY . /go/src/app
WORKDIR /go/src/app

# Get dependencies and compile App
RUN go-wrapper download && \
          go-wrapper install

# Expose ports and commands
EXPOSE 8601
CMD ["go-wrapper", "run"]