FROM hub.iot.chinamobile.com/offline/librdkafka-golang1141-alpine:latest
WORKDIR /TESTRUN/
COPY . .
ENV PHASE 1
RUN go build -mod=vendor -o build ./cmd/build/build.go && \
    go build -mod=vendor -o grpctest ./cmd/grpc/grpc.go && \
    go build -mod=vendor -o pipeline ./cmd/pipeline/pipeline.go && \
    go build -mod=vendor -o artifact ./cmd/artifact/artifact.go && \
    go build -mod=vendor -o pods ./cmd/pods/pods.go && \
    go build -mod=vendor -o vm ./cmd/vm/vm.go && \
    go build -mod=vendor -o probe ./cmd/probe/probe.go &&\
    go build -mod=vendor -o test ./cmd/test/test.go && \
    go build -mod=vendor -o common ./cmd/common/common.go && \ 
    go build -mod=vendor -o concurrent ./cmd/concurrent/concurrent.go && \
    go build -mod=vendor -o concurrent2 ./cmd/concurrent2/concurrent2.go 

ENV PHASE 2
RUN chmod +x run_test.sh
CMD [ "/TESTRUN/run_test.sh" ]