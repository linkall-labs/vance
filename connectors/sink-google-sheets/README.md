---
title: Google Sheet
---

# Google Sheet Sink

## Introduction

The Google Sheet Sink is a [Vanus Connector][vc] which aims to send incoming CloudEvents in a JSON format to a Google
Sheet.

For example, if the incoming CloudEvent looks like::

```json
{
  "id": "88767821-92c2-477d-9a6f-bfdfbed19c6a",
  "source": "quickstart",
  "specversion": "1.0",
  "type": "quickstart",
  "time": "2022-07-08T03:17:03.139Z",
  "datacontenttype": "application/json",
  "data": {
    "id": "1",
    "name": "Ehis",
    "email": "ehis@gmail.com",
    "description": "Developer"
  }
}
```

The Google Sheet Sink will extract data field and write it to a Google Sheet.

## Pre-requisite

- A Google Sheet
- Service account on the google cloud platform for server authentication - Ensure you give Service Account Editor Access

**Note:** It’s necessary to share the spreadsheet with client_email of the service account to access it. Otherwise, you
will get 403 forbidden. You can find client_email in the downloaded key’s json file

### Create the config file

```shell
cat << EOF > config.yml
credentials: |-
  {
  "type": "service_account",
  "project_id": "user-auth-123456",
  "private_key_id": "2ce83c7aa4fc3ad0a613862asladfsafas",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIEvAIBADANBgkqhkiG9w0BAQEFAASCBKYwggSiAgEAAoIBAQCypWa6Oz/3d+XX\n06/eI3c4AaZs3x3ZCeXStGqYcOXx2mQEG8ro2UILW2W63p+TfqJtt0eh4yM+zcLu\lasdflajLSAJLDSFJJKSljljalfjslfS\np92aEGnYzSgsSv2J6dYXv/JfleOBxKrHgfsMO59lqhCAa/NFOqau+L+yOJsQ4atG\nw1/onj4hqNDpK3uDxhYI3RCWow12b6IVV84UTmdRomX9BAlSgz0IK7EllhQuL+VP\nXud4TKgmLboshL5V62B49MYDf1WgAi8YjvqHBepbl0d0DJvI4Uxsx/K0ORXzFH4C\n6oA1FvxhAgMBAAECggEABMiCdWKQPiVG8W4w04sSDl0W8mvP7geCkI91mjLNqVne\nWzDOUEk+6DRwhx4eWjHmEfd6EsbU4wHZ4g1BAEoD1urWL0lf2GKO7JgNY0S1ReXT\nW9TGohJ/jMBtgPziVHtE3EcQKrQO5ATCCo7cU1t0phIqLPNuEfMBoT2ptO77ujX5\nFawpt8ctPmaHP7g1TmKuRRL5hDavEnopx1BIxstX1oTZtPY55/jEdQjrYipTfi1M\nQfeaZM6xKFtfBgAjHQiJiOUA65Jyqp38qBJVfkpCLZbNFtL4X6lNKWx9ZCFYcSbY\nIUwuM6acq92UYl3aikg7RJRxD7UtlrvWnlHprdCkRQKBgQDr/5cOMkWDBPJv1uiP\nXMy0dI9Qu881XMn78X24td9KP/Qx0c6XpsTf1nYykySBI5KZkw4fh8JSJvF63jZH\nUiObz3OBH4fvonpK9Az42KiVep8gctCQtXgLEZd2QLcy0r7OWojE21JP9T+BP0OL\nDGR+OWYW5t0GD7VaqnATSw+AfQKBgQDByXMqbOyWn0Oo8RQSmOKP+IHIqEIvJ1j4\nPq+8/01lCg2UzAVEHyPxnv+NgbHg1HGq/c3Ez2zUtnxXt6w5GHW43ZsSuaosag82\nVzyMv/3faYIsgjh4qvcruCpQSloNNeW2STG4qvZ7qCOVDJg0tvoeEP6iC2aYhtoJ\nBqEoyXL0tQKBgCOCaLMtI0JsiyIC3yk7GF4Kr8nBCJOJ66ZqFrWlP/zHFLIuVHyD\nDlpzxYMkmriHprZO5zAdWELOM0V+jAI9PLhkBYgnO2f2NZpzkEQdLXiYY7sZK4Kq\nm25m7jhP0oDmLumTu8KLEZ6QU0baQwp4CeLoNhE6GYWg9XO383cjsyhtAoGAcRxR\nsWjEq6IojvqwWa6NR7Woo2O6xeU0pCmK0ElAdoJorPps9HcstsK0rXcPSYkXE9Ry\n/7aG8p3VdCnMR8NEK3SGKGbgsm3xlSlUOV9zIq1mAu67YYuBHC6x3A2aBG36N+z/\nLaf0mPbqVfx09wf6dAQ9bH41E0BbEbuh47m59KUCgYAljW2YcQS+ReCzeeBcpcxt\ng/uRVgfBDY/eYnlMRRFMGL5Jg2BrQDxkb4VU+BzNlbPK41UmUAfa0/OB8uT80bG3\nJf4aIWyZU+AkAXG7MIwR+ZMD1RVSmPIo2X44nddrpMX2he+AsuM2+Xwbr+q18nWv\n3N1Vb31GxK/iD8Pw3ItPIQ==\n-----END PRIVATE KEY-----\n",
  "client_email": "sheetauth@user-auth-123456.iam.gserviceaccount.com",
  "client_id": "1000161824343520234567",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/sheetauth%40user-auth-123456.iam.gserviceaccount.com"
  }
sheet_id : https://docs.google.com/spreadsheets/d/1BxiMVs0XRA5nFMdKvBdBZjgmUUqptlbs74OgvE2upms/edit#gid=0
sheet_name: vanus

EOF
```

| Name        | Required |   Default    | Description                                   |
|:------------|:--------:|:------------:|-----------------------------------------------|
| port        |    NO    |     8080     | the port which Google Sheets Sink listens on  |
| credentials |   YES    |              | Google [Service Account][sa] credentials JSON |
| sheet_url   |   YES    |              | Google sheet url                              |

[vc]: https://docs.vanus.ai/introduction/concepts#vanus-connect
[sa]: https://developers.google.com/workspace/guides/create-credentials?hl=zh-cn#service-account
