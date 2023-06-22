#In fact, your Dockerfile is leveraging multi-stage builds, 
# a powerful feature that allows you to separate your build into multiple stages, 
#with each stage beginning with a FROM statement.

#Here's how your Dockerfile is working:

#    The FROM golang:latest AS backend_builder statement is creating a stage where you're building your Go app. 
#    This stage uses a base image with Go installed.

#    The FROM node:latest as frontend_builder statement is creating another stage where you're building your Angular frontend. 
#   This stage uses a base image with Node.js installed.

#    The FROM alpine:latest statement is creating the final stage where you're bundling the built Go app and the built Angular app into a single image. 
#    This stage uses Alpine Linux, a very lightweight Linux distribution, as its base image.

#At the end, only the final stage (the one with FROM alpine:latest) will be kept, and the intermediate images created by the other stages will be discarded. 
#This allows you to build your app using necessary tools and libraries, and then bundle only the built app into a small, lightweight image.

#So, while the FROM alpine:latest statement is at the end of your Dockerfile, it's not installed after everything else. 
#Rather, it's the base image for the final Docker image that is produced from your Dockerfile. 
#All the installations and builds happen in separate, intermediate stages, and only the results (the built Go app and Angular app) are copied into the final image.


# Start from a Debian image with the latest version of Go installed
# and a workspace (GOPATH) configured at /go.
FROM golang:latest AS backend_builder

# Copy the local package files to the container's workspace.
ADD . /go/src/app

# Build the Go app
WORKDIR /go/src/app
RUN go get -d -v ./...
RUN go install -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Start with a base image containing Javascript runtime
RUN mkdir /hermex
RUN mkdir /hermex/rss-ui
RUN mkdir /hermex/rss-ui/dist
RUN mkdir /hermex/rss-ui/dist/rss-index

# Final stage
FROM alpine:latest

# Install ca-certificates so that HTTPS works consistently
RUN apk update && apk add ca-certificates && rm -rf /var/cache/apk/*


WORKDIR /hermex
COPY --from=backend_builder /go/src/app/main .
COPY rss-ui/dist/rss-index rss-ui/dist/rss-index 


RUN chmod -R 777 /hermex

# Set the command to run your app using CMD which defines your runtime
# Here we are using the "main" binary that was created in the previous step as our entrypoint
CMD ["/hermex/main"]

# Expose port 8080 to the outside world
EXPOSE 8080
