# Installation

this package requires the following docker images:

* rabbitmq:3-management
* mailcatcher

```
docker run -d --name styx-rabbit -p 5672:5672 -p 15672:15672 rabbitmq:3-management
docker run -d --name styx-mailcatcher -p 1025:1025 -p 1080:1080 zolweb/docker-mailcatcher
```

## Access to RabbitMQ
Url: `http://localhost:15672` - guest/guest

## Access to Mailcatcher
Url: `http://localhost:1080`

# API

Submit an E-Mail:

POST: localhost:9999/api/mail

Body:
```
{
	"data": {
		"type": "mail",
		"attributes": {
			"context": "foo",
			"subject": "test mail",
			"clients": [
					{
						"name": "Johannes Pichler",
						"email": "johannes.pichler@karriere.at",
						"type": "to"
					},
					{
						"name": "Johannes Pichler",
						"email": "johannes.pichler@jopic.at",
						"type": "to"
					},
					{
						"name": "Johannes Pichler",
						"email": "johannes.pichler@karriere.at",
						"type": "from"
					}
				],
			"body": {
				"html": "html here",
				"plain": "plain text here"
			},
			"priority": 1,
			"attachments": []
		}
	}
}
```