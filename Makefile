SRC = server.go routing.go logging.go
TRG = gogate

build:
	go build -o ${TRG} ${SRC}

run:
	go run ${SRC}

clean:
	rm -f ${TRG}

