FROM dtcooper/raspberrypi-os:bookworm As builder

RUN apt-get update && apt-get upgrade -y && \
    apt-get install -y curl gcc git locales locales-all &&\
    locale-gen ja_JP.UTF-8
RUN curl -OL https://go.dev/dl/go1.22.1.linux-arm64.tar.gz &&\
tar -C /usr/local -xzf go1.22.1.linux-arm64.tar.gz &&\
rm -rf go1.22.1.linux-arm64.tar.gz
ENV PATH $PATH:/usr/local/go/bin

WORKDIR /app
ADD ./main .
RUN go build -o app

FROM dtcooper/raspberrypi-os:bookworm

RUN apt-get update && \
    apt-get upgrade -y && \
    apt-get install -y locales locales-all curl && \
    apt-get clean &&\
    rm -rf /var/lib/apt/lists/*
WORKDIR /app
COPY --from=builder /app/app /app
CMD ["./app"]