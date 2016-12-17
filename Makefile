all:
	go get ./...
	go build -buildmode=c-shared -o out_kafka.so .

clean:
	rm -rf *.so *.h *~
