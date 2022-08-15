#FROM golang:1.14 as go14
FROM golang:1.16
WORKDIR /app
COPY . .

# Replace shell with bash so we can source files
RUN rm /bin/sh && ln -s /bin/bash /bin/sh

RUN apt-get clean && apt-get update
RUN apt-get install ca-certificates libgnutls30 -y
RUN cd build && go build -o bin ./manifest
RUN cd build && go build -o bin ./pluginctl

# Install nvm and node
#ENV NVM_DIR /usr/local/nvm
#ENV NODE_VERSION 14.19.3
#RUN mkdir -p ${NVM_DIR}
#RUN curl https://raw.githubusercontent.com/nvm-sh/nvm/v0.39.1/install.sh | bash \
#    && . $NVM_DIR/nvm.sh \
#    && nvm install $NODE_VERSION \
#    && nvm alias default $NODE_VERSION \
#    && nvm use default
#ENV NODE_PATH $NVM_DIR/v$NODE_VERSION/lib/node_modules
#ENV PATH      $NVM_DIR/versions/node/v$NODE_VERSION/bin:$PATH

RUN curl -fsSL https://deb.nodesource.com/setup_16.x | bash \
    && apt-get install -y nodejs

#RUN apt install -y python

#RUN npm install node-sass@4.14.1 --save-dev

#FROM golang:1.16
#WORKDIR /app
#COPY --from=go14 . .
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.31.0
RUN make

ENTRYPOINT ["tail", "-f", "/dev/null"]
