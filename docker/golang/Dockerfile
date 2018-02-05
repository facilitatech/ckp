FROM golang:1.9

# app specific env
ENV BUILDPATH /go/src/
ENV TOOLS /var/exec
ENV GOBINARIES /go/bin

# Install reflex
WORKDIR $GOBINARIES
RUN apt-get update && apt-get upgrade -y && apt-get install vim -y

RUN go get github.com/cespare/reflex
RUN go get github.com/agtorre/gocolorize
RUN cd /go/src/ && mv github.com /usr/local/go/src/

# Make directories and add files as needed
RUN mkdir -p $TOOLS
ADD build.sh $TOOLS
ADD reflex.conf $TOOLS
RUN chmod +x $TOOLS/build.sh

# Execute reflex.
WORKDIR $BUILDPATH
CMD ["reflex","-c","/var/exec/reflex.conf"]