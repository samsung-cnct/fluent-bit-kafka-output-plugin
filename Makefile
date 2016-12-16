all:
	go build -buildmode=c-shared -o out_kafka.so .
