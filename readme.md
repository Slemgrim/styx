# Installation

To but up a basic environment use docker compose
```
cp config-sample.json config.json
docker-compose up
```

## Access to RabbitMQ
Url: `http://localhost:15672` - guest/guest

## Access to Mailcatcher
Url: `http://localhost:1080`

# API

Submit an E-Mail:

POST: localhost:9999/mails

Body:
```json
{
	"data": {
		"type": "mail",
		"attributes": {
			"subject": "test mail",
			"to": [
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				}
			],
			"from": {
				"name": "Rick Sanchez",
				"address": "rick@sanchez.com"
			},

			"body": {
				"html": "<h1>I'm the body</h1>"
			}
		}
	}
}
```

Submit a more complex E-Mail:

POST: localhost:9999/mails

Body:
```json
{
	"data": {
		"type": "mail",
		"attributes": {
			"subject": "test mail",
			"to": [
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				},
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				}
			],
			"cc": [
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				},
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				}
			],
			"bcc": [
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				},
				{
					"name": "Rick Sanchez",
					"address": "rick@sanchez.com"
				}
			],
			"from": {
				"name": "Rick Sanchez",
				"address": "rick@sanchez.com"
			},
			"reply-to": {
				"name": "Rick Sanchez",
				"address": "rick@sanchez.com"
			},
			"return-path": {
				"name": "Rick Sanchez",
				"address": "rick@sanchez.com"
			},
			"body": {
				"html": "html here",
				"plain": "plain text here"
			},
			"headers": [
				{
					"name": "foo",
					"value": ["foo", "bar"]
				},
				{
					"name": "foo",
					"value": ["foo"]
				}
			]
		}
	}
}
```

## Attachments

First you have to create an attachment resource

POST: localhost:9999/attachments

Body:
```json
{
	"data": {
		"type": "attachment",
		"attributes": {
			"size": 4,
			"file-name": "test.txt",
			"mime-type": "text/plain; charset=utf-8",
			"hash": "3da541559918a808c2402bba5012f6c60b27661c"
		}
	}
}
```

Make sure that the size, mime-type and hash (sha1) match your file.
You will receive an upload path with the response.
Use this upload path to upload your file

You can use references to attachments within your mail

Body:
```json
{
	"data": {
		"type": "mail",
		"attributes": {
			...
		},
		"relationships": {
            "attachments": {
                "data": [
                    {
                        "type": "attachment",
                        "id": "c2319772-172d-469b-82bb-28920b6c786c"
                    }
                ]
            }
        }
	}
}
```
