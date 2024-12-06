FROM golang:latest

# Create the working directory inside the container
WORKDIR /ccommits

# Copy source files and build the application
COPY . .
RUN go build -o ccommits
RUN mv ccommits/conventional-commits-cli /usr/bin/

# Installing the go debugger delve
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN dlv version

# Sets the TERM environment to match requirements by ccommits
ENV TERM=xterm-256color

# Run the application
CMD [ "conventional-commits-cli" ]