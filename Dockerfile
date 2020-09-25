# Note: We don't use Alpine and its packaged Rust/Cargo because they're too often out of date,
# preventing them from being used to build Substrate/Polkadot.

# First Phase - Load Subkey
FROM parity/subkey:2.0.0 as subkey
RUN subkey --version

## Second Phase - Build context for tests
FROM parity/substrate:v2.0.0-rc6

USER root

COPY --from=subkey /usr/local/bin/subkey  /usr/local/bin/subkey

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
		git-core \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.14.9

RUN set -eux; \
	\
# this "case" statement is generated via "update.sh"
	dpkgArch="$(dpkg --print-architecture)"; \
	case "${dpkgArch##*-}" in \
		amd64) goRelArch='linux-amd64'; goRelSha256='f0d26ff572c72c9823ae752d3c81819a81a60c753201f51f89637482531c110a' ;; \
		armhf) goRelArch='linux-armv6l'; goRelSha256='e85dc09608dc9fc245ebc5daea0826898ac0eb0d48ed24e2300427850876c442' ;; \
		arm64) goRelArch='linux-arm64'; goRelSha256='65e6cef5c474a3514e754f6a7987c49388bb85a7b370370c1318087ac35427fa' ;; \
		i386) goRelArch='linux-386'; goRelSha256='14982ef997ec323023a11cffe1a4afc3aacd1b5edebf70a00e17b67f888d8cdb' ;; \
		ppc64el) goRelArch='linux-ppc64le'; goRelSha256='5880a37faf93b2396edc3ff231e0f8df14d0520505cc13d01116e24d7d1d0147' ;; \
		s390x) goRelArch='linux-s390x'; goRelSha256='381fc24aff153c4affcb00f4547683212157af29b8f9e3de5952d78ac35f5a0f' ;; \
		*) goRelArch='src'; goRelSha256='c687c848cc09bcabf2b5e534c3fc4259abebbfc9014dd05a1a2dc6106f404554'; \
			echo >&2; echo >&2 "warning: current architecture ($dpkgArch) does not have a corresponding Go binary release; will be building from source"; echo >&2 ;; \
	esac; \
	\
	url="https://golang.org/dl/go${GOLANG_VERSION}.${goRelArch}.tar.gz"; \
	curl -o go.tgz -L "$url"; \
	echo "${goRelSha256} *go.tgz" | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	if [ "$goRelArch" = 'src' ]; then \
		echo >&2; \
		echo >&2 'error: UNIMPLEMENTED'; \
		echo >&2 'TODO install golang-any from jessie-backports for GOROOT_BOOTSTRAP (and uninstall after build)'; \
		echo >&2; \
	exit 1; \
	fi; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH
WORKDIR $GOPATH

RUN mkdir -p $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
WORKDIR $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
COPY . .

# Ensuring Subkey is available
RUN subkey --version

RUN make install

# Reset parent entrypoint
ENTRYPOINT []
CMD ["make", "test-cover"]
