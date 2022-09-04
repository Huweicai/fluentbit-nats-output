all: image

so:
	go build -buildmode=c-shared -o nats-output.so .

image: so
	docker build -t fluent/fluent-bit:1.9-nats .

clean:
	rm -rf *.so *.h *~
