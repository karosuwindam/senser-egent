TAG = 0.0.1
DOCKER = docker
NAME = sennser-egent

TEMPLATE = ./Dockerfile_tmp
TARGET = ./Dockerfile
TARGET_FILE = ./main
GO_VERSION = 1.22.1
ARCH = arm64
OPT = "--privileged"

TRACER_ON = true
TRACER_GRPC_URL = otel-grpc.bookserver.home:4317


all: build
create:
	@echo "--- ${NAME} ${TAG} create ---"
	@echo "--- create Dockerfile ---"
	@cat ${TEMPLATE} | sed s/TAG/${TAG}/ | sed s/ARCH/${ARCH}/ | sed s/GO_VERSION/${GO_VERSION}/ > ${TARGET}
build: create
	@echo "--- build Dockerfile --"
	${DOCKER} build -t ${NAME}:${TAG} -f ${TARGET} .
rm: 
	${DOCKER} rmi ${NAME}:${TAG}
run:
	${DOCKER} run --rm --name=${NAME} ${OPT} -e TRACER_ON:${TRACER_ON} -e TRACER_GRPC_URL:${TRACER_GRPC_URL} -p 18080:8080 ${NAME}:${TAG}
push:
	${DOCKER} push ${NAME}:${TAG}
