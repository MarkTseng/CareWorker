PROG_BIN=CareWorker

all:
	cd server && go build -o ../${PROG_BIN}
clean:
	rm -rf ${PROG_BIN}
