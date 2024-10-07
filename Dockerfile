FROM golang:1.22-alpine


ARG ENV


#RUN if [ "$ENV_ARG" == "production" ] ; then  \
#      cp .docker/dev/start-container.sh /tmp/start-container.sh ;  \
#    else cp .docker/prod/start-container.sh /tmp/start-container.sh ; fi



# Setup the Go app

RUN mkdir /golangApp
COPY . /golangApp
WORKDIR /golangApp

# this file should be for the develoment only because of the CompileDaemon
RUN apk add --update --no-cache vim bash git nano htop curl
RUN go clean -modcache

RUN go install github.com/githubnemo/CompileDaemon@latest






COPY .docker/${ENV}/start-container.sh /tmp/start-container.sh

RUN chmod +x /tmp/start-container.sh

ENTRYPOINT ["/tmp/start-container.sh"]

# Expose the port the app will run on
EXPOSE 8080

