FROM fluent/fluent-bit:1.9-debug
ADD ./nats-output.so /fluent-bit/bin/nats-output.so