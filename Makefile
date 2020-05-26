PROG_BIN=CareWorker

all:
	cd server && go build -ldflags="-s -w" -o ../${PROG_BIN}
	cd ../
	#goupx ${PROG_BIN}
	./${PROG_BIN}
clean:
	rm -rf ${PROG_BIN}
