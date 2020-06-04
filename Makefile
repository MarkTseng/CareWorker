PROG_BIN=CareWorker

all: build
	@echo "build Careworker web server done"
	@echo "For upx binary, Please make goupx"

goupx: 
	goupx ${PROG_BIN}

run:
	./${PROG_BIN}

build: 
	@echo "build Careworker web server start"
	cd server && go build -ldflags="-s -w" -o ../${PROG_BIN}

clean:
	rm -rf ${PROG_BIN}
