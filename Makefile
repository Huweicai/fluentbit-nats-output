all: image

so:
	go build -buildmode=c-shared -o ./build/output/nats-output.so .

image: so
	docker build -f ./build/Dockerfile -t fluent/fluent-bit:1.9-nats .

clean:
	rm -rf *.so *.h *~
