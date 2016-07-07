rel:
	go build

debug:
	go build -tags "debug"

run:
	sudo ./ok-mysql --addr="127.0.0.1:3306"
