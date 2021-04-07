# Note: We don't use Alpine and its packaged Rust/Cargo because they're too often out of date,
# preventing them from being used to build Substrate/Polkadot.

## First Phase - Build context for tests
FROM parity/substrate:v3.0.0

USER root

# gcc for cgo
RUN apt-get update && apt install -y --no-install-recommends software-properties-common dirmngr gnupg2
RUN apt-get update && apt-key adv --keyserver keyserver.ubuntu.com --recv-keys F6BC817356A3D45E
RUN add-apt-repository ppa:longsleep/golang-backports && apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
		git-core \
		golang-go ## install the latest go \
	&& rm -rf /var/lib/apt/lists/*

RUN go version
ENV GOPATH /go
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR $GOPATH

RUN mkdir -p $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
WORKDIR $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
COPY . .

RUN make install

# Reset parent entrypoint
ENTRYPOINT []
CMD ["make", "test-cover"]
