# Feishu Sink 

## Introduction

The Feishu Sink is a [Vance Connector](../README.md) which aims to handle incoming CloudEvents in a way that extracts the `data` part of the
original event and deliver these extracted `data` to  Feishu APIs.

For example, if the incoming CloudEvent looks like:

```http
{
  "id" : "42d5b039-daef-4071-8584-e61df8fc1354",
  "source" : "vance-http-source",
  "specversion" : "V1",
  "type" : "http",
  "datacontenttype" : "application/json",
  "time" : "2022-05-17T18:44:02.681+08:00",
  "vancefeishusinkservice": "bot",
  "data" : {
    ...
  }
}
```

### Supported Feishu Service

- Bot: pushing a message to Group Chat

## Quick Start

in this section, we show how to use Feishu Sink push a message to your group chat.

### Add a bot to your group chat

Go to your target group, click Chat Settings > Group Bots > Add Bot, and select Custom Bot to add the bot to the group chat.

Enter a name and description for your bot, or set up an avatar for the bot, and then click "Add".

![add-a-bot](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/add-a-bot.gif)

You will get the webhook address of the bot in the following format:

```
https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx
```

> ⚠️ Please keep this webhook address properly. Do not publish it on GitHub, blogs, and other publicly accessible sites to avoid it being maliciously called to send spam messages.

![bot-config](https://github.com/linkall-labs/vance-docs/raw/main/resources/connectors/sink-feishu-bot/feishu-config.png)

>  ⚠️ You must set your signature verification to make sure push messages work.

### Create Config file
```shell
cat << EOF > config.yml
# change the webhook and bot_signature to your.
secret:
  bot_signature: "xxxxxx"
enable: ["bot"]
bot:
  webhook: "https://open.feishu.cn/open-apis/bot/v2/hook/xxxxxxxxxxxxxxxxx"
EOF
```

### Start Using Docker

```shell
docker run -d --rm \
  --network host \
  -v ${PWD}:/vance/config \
  --name sink-feishu public.ecr.aws/vanus/connector/sink-feishu:dev
```

### Test

```json
curl --location --request POST '127.0.0.1:8080' \
--header 'Content-Type: application/cloudevents+json' \
--data-raw '{
  "id": "53d1c340-551a-11ed-96c7-8b504d95037c",
  "source": "sink-feishu-quickstart",
  "specversion": "1.0",
  "type": "hello",
  "datacontenttype": "application/json",
  "time": "2022-10-26T10:38:29.345Z",
  "vancefeishusinkservice": "bot",
  "data": "Hello Feishu!"
}'
```

now, you cloud see a notice in your chat group.
[received-notification](received-message.png)


## Configuration

The default path is `/vance/config/config.yml`. if you want to change the default path, you can set env `CONNECTOR_CONFIG` to
tell Feishu Sink.

| Name     | Required | Default | Description                                     |
|:---------|:--------:|:-------:|-------------------------------------------------|
| secret.bot_signature | **YES**  |    -    |  Feishu Bot signature.| "https://open.feishu.cn/open-apis/bot/v2/hook/......                       |
| enable     | **YES**  |    -    | which services you want Feishu Sink are enabled |
| bot.webhook     | **YES**  |    -    | HTTP endpoint of Feishu Bot, looks like https://open.feishu.cn/open-apis/bot/v2/hook/...... |

### Secret(Optional)
If you want separate secret information to an independent file, you could create a file like:

```shell
cat << EOF > secret.yml
bot_signature: "xxxxxx"
EOF
```
then mount it into your container. The default path of it is `/vance/config/secret.yml`. if you want to change the default path,
you can set env `CONNECTOR_SECRET` to tell Feishu Sink.

This feature is very useful when you want to use [Secrets](https://kubernetes.io/docs/concepts/configuration/secret/) in Kubernetes