# Info
FROM golang:latest
LABEL MAINTAINER="Alexandru-Paul Copil <alexandru.p.copil@gmail.com>"

# Open ports to the network
EXPOSE 2000 9000

# Set project path
ENV CRYPTOR_PATH=/go/src/github.com/thee-engineer/cryptor

# Copy project components
RUN mkdir -p ${CRYPTOR_PATH}
ADD . ${CRYPTOR_PATH}

# Set the workdir
WORKDIR ${CRYPTOR_PATH}