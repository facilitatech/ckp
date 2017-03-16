FROM golang:1.7
# app specific env
ENV BUILDPATH /go/src/app
ENV TOOLS /var/exec
ENV GOBINARIES /go/bin

# Install reflex
WORKDIR $GOBINARIES
RUN go get github.com/cespare/reflex

ENV PORT 6060
# docker/app port
EXPOSE $PORT

# Make directories and add files as needed
RUN mkdir -p $TOOLS
ADD build.sh $TOOLS
ADD reflex.conf $TOOLS
RUN chmod +x $TOOLS/build.sh

# Execute reflex.
WORKDIR $BUILDPATH
CMD ["reflex","-c","/var/exec/reflex.conf"]
