FROM golang:latest

# Create the working directory inside the container
WORKDIR /ccommits

# Copy source files and build the application
COPY . .
RUN go build -o ccommits
RUN mv ccommits/ccommits /usr/bin/

# Sets the TERM environment to match requirements by ccommits
ENV TERM=xterm-256color

# Run the application
CMD [ "ccommits" ]