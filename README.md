# FluentBit NATS Output Plugin

## Features

- Support credentials
- Auto sniffer plain and tls protocol to connect with
- Support compression
- Reconnecting while server not available
- Support multi output instance

## Config parameters

| Parameter Name | Meaning                                    |
| -------------- | ------------------------------------------ |
| Name           | Plugin name, must be **nats-output**       |
| URL            | NATS connection URL                        |
| Subject        | The NATS subject used to publish data      |
| TimeoutSeconds | NATS connection dial timeout in seconds    |
| Compression    | Whether to enable transmission compression |

Config demo:

```
[OUTPUT]
    Name  nats-output
    Match file.*
    URL nats://admin:PassWord@192.168.2.57:4222
    Subject fluentbit.test
    TimeoutSeconds 5
    Compression true
```