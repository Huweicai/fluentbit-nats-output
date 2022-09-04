FROM fluent/fluent-bit:1.9-debug
ADD ./build/output/nats-output.so /fluent-bit/bin/nats-output.so
ADD ./build/plugins.conf /fluent-bit/etc/plugins.conf