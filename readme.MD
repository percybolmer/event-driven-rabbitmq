# An repository showcasing Event Driven Architecture (EDA)
This repository is an example of how to use EDA.

You can read about this repository and how it works in 

https://programmingpercy.tech/blog/event-driven-architecture-using-rabbitmq/

## Running it

Generate the certificates
```bash
git clone https://github.com/rabbitmq/tls-gen tls-gen
cd tls-gen/basic
make PASSWORD=
make verify
```

Replace the information in `rabbitmq.conf` so that it points towards
your generated certificate names

Start Docker and mount files 

```bash
sudo docker run -d --name rabbitmq -v "$(pwd)"/rabbitmq_definitions.json:/etc/rabbitmq/rabbitmq_definitions.json:ro -v "$(pwd)"/rabbitmq.conf:/etc/rabbitmq/rabbitmq.conf:ro -v "$(pwd)"/tls-gen/basic/result:/certs -p 5671:5671 -p 15672:15672 rabbitmq:3.11-management
```