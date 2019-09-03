# Note: We don't use Alpine and its packaged Rust/Cargo because they're too often out of date,
# preventing them from being used to build Substrate/Polkadot.

FROM cd544298d051

# gcc for cgo
RUN apt-get update && apt-get install -y --no-install-recommends \
		g++ \
		gcc \
		libc6-dev \
		make \
		pkg-config \
	&& rm -rf /var/lib/apt/lists/*

ENV GOLANG_VERSION 1.12.9

RUN set -eux; \
	\
# this "case" statement is generated via "update.sh"
	dpkgArch="$(dpkg --print-architecture)"; \
	case "${dpkgArch##*-}" in \
		amd64) goRelArch='linux-amd64'; goRelSha256='ac2a6efcc1f5ec8bdc0db0a988bb1d301d64b6d61b7e8d9e42f662fbb75a2b9b' ;; \
		armhf) goRelArch='linux-armv6l'; goRelSha256='0d9be0efa9cd296d6f8ab47de45356ba45cb82102bc5df2614f7af52e3fb5842' ;; \
		arm64) goRelArch='linux-arm64'; goRelSha256='3606dc6ce8b4a5faad81d7365714a86b3162df041a32f44568418c9efbd7f646' ;; \
		i386) goRelArch='linux-386'; goRelSha256='c40824a3e6c948b8ecad8fe9095b620c488b3d8d6694bdd48084a4798db4799a' ;; \
		ppc64el) goRelArch='linux-ppc64le'; goRelSha256='2e74c071c6a68446c9b00c1717ceeb59a826025b9202b3b0efed4f128e868b30' ;; \
		s390x) goRelArch='linux-s390x'; goRelSha256='2aac6de8e83b253b8413781a2f9a0733384d859cff1b89a2ad0d13814541c336' ;; \
		*) goRelArch='src'; goRelSha256='ab0e56ed9c4732a653ed22e232652709afbf573e710f56a07f7fdeca578d62fc'; \
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

# For debugging
# compile Delve
# RUN apt-get update && apt-get install -y git
# RUN go get github.com/derekparker/delve/cmd/dlv

# Port 40000 belongs to Delve
# EXPOSE 40000

# Allow delve to run on Alpine based containers.
# RUN apk add --no-cache libc6-compat

# FROM phusion/baseimage:0.10.2
# LABEL maintainer="chevdor@gmail.com"
# LABEL description="This is the 2nd stage: a very small image where we copy the subkey binary."
# ARG PROFILE=release

# RUN mv /usr/share/ca* /tmp && \
# 	rm -rf /usr/share/*  && \
# 	mv /tmp/ca-certificates /usr/share/ && \
# 	mkdir -p /root/.local/share/Polkadot && \
# 	ln -s /root/.local/share/Polkadot /data && \
# 	useradd -m -u 1000 -U -s /bin/sh -d /substrate substrate

# COPY --from=builder /substrate/target/$PROFILE/substrate /usr/local/bin

# # checks
# RUN ldd /usr/local/bin/substrate && \
# 	/usr/local/bin/substrate --version

# # Shrinking
# RUN rm -rf /usr/lib/python* && \
# 	rm -rf /usr/bin /usr/sbin /usr/share/man

# USER substrate
# EXPOSE 30333 9933 9944
# VOLUME ["/data"]

# CMD ["/usr/local/bin/substrate"]

RUN mkdir -p $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
WORKDIR $GOPATH/src/github.com/centrifuge/go-substrate-rpc-client
COPY . .

RUN go get -d -v ./...

CMD ["go", "run", "test/main.go"]
# CMD ["dlv", "debug", "--listen=:40000", "--headless=true", "--api-version=2", "--log", "github.com/centrifuge/go-substrate-rpc-client/test"]
