FROM public.ecr.aws/lambda/provided:al2 as build

# install build deps
RUN yum install -y \
    git \
    gcc \
    gcc-c++ \
    make \
    wget \
    tar \
    gzip \
    pkgconfig \
    glib2-devel \
    expat-devel && \
    # install golang
    GOLANG_VERSION="1.17.3" && \
    wget "https://dl.google.com/go/go$GOLANG_VERSION.linux-amd64.tar.gz" && \
    rm -rf /usr/local/go && \
    tar -C /usr/local -xzf "go$GOLANG_VERSION.linux-amd64.tar.gz" && \
    rm -rf "go$GOLANG_VERSION.linux-amd64.tar.gz"

# go build
# 実行アーキテクチャ: x86_64 (amd64)
ENV PATH $PATH:/usr/local/go/bin
ADD . .
WORKDIR app
RUN go build -o /main main.go

# マルチステージビルド
FROM public.ecr.aws/lambda/provided:al2
COPY --from=build /main /main
ENTRYPOINT [ "/main" ]
