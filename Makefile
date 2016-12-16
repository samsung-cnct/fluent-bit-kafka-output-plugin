all:
	go build -buildmode=c-shared -o out_kafka.so .

clean:
	rm -rf *.so *.h *~
