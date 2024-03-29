# teadiller


To run golang inside of docker container. You need to run the following command in project folder.
```
docker run --rm -it -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" golang bash
```
OR
compile and run
```
docker run --rm -t -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" 
--env TELEGRAMM_TOKEN="<your_token>"
--env TELEGRAM_BOT_NAME="<your_botname>" 
--env TELEGRAM_BOT_WEB_PORT="<your_web_port>" 
--env TELEGRAM_BOT_DB_PATH="<path to xlsx>"
--env PAYPAL_SANDBOX="<true/false, true for running in sandbox mode>"
--env PAYPAL_USERNAME="<your paypal username>"
--env PAYPAL_PASSWORD="<your paypal password>"
--env PAYPAL_SIGNATURE="<your paypal signature>"
--env PAYPAL_OK_URL="<redirect to this url if everything is ok...>"
--env PAYPAL_CANCEL_URL="<...otherwise, redirect to this>"
--env GEOCODER_API_KEY="<get api key at http://mapquestapi.com>"
-p $TELEGRAM_BOT_WEB_PORT:$TELEGRAM_BOT_WEB_PORT golang go run main.go
```
Run detached 
```
docker run -d -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go"
--env TELEGRAMM_TOKEN="<your_token>"
--env TELEGRAM_BOT_NAME="<your_botname>"
--env TELEGRAM_BOT_WEB_PORT="<web_port>"
--env TELEGRAM_BOT_DB_PATH="<path to xlsx>"
--env PAYPAL_SANDBOX="<true/false, true for running in sandbox mode>"
--env PAYPAL_USERNAME="<your paypal username>"
--env PAYPAL_PASSWORD="<your paypal password>"
--env PAYPAL_SIGNATURE="<your paypal signature>"
--env PAYPAL_OK_URL="<redirect to this url if everything is ok...>"
--env PAYPAL_CANCEL_URL="<...otherwise, redirect to this>"
--env GEOCODER_API_KEY="<get api key at http://mapquestapi.com>"
-p $TELEGRAM_BOT_WEB_PORT:$TELEGRAM_BOT_WEB_PORT golang go run main.go
```

## Configurations
| Env variable name | Value |
|-------------------|-------|
| TELEGRAM_TOKEN   |\<Put telegram token here\>|
| TELEGRAM_BOT_NAME|\<Put telegram bot name\>|
| TELEGRAM_BOT_WEB_PORT      |\<Put web port\>|
| TELEGRAM_BOT_DB_PATH |\<Put path to xlsx file\>|

## Instagram api examples
```
msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{[][]string{{"black tea", "green tea", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8"}}, false, false, false}
```

## Excel file format
There should be at least one sheet - 'items', which is filled by user.

This sheet must contain the next columns (order is important):

| Id | Name | Description | Photo | Tags | Item Type | Properties | Count |
|----|------|-------------|-------|------|-----------|------------|-------|

Id - numerical identificator of item

Name - displayed name of item

Description - some description

Photo - path to image of item

Tags - some tags, divided by comma (e.g.: tea,green,indian)

Item Type - unit/measure of item (e.g.: gramm, liter, unit and etc.)

Properties - additional properties, defined in the next format: 'prop1=val1,prop2=val2,...'

Count - amount of available ```Item Type```

There's also one more sheet - 'orders' with the next columns:

| Id | Item Id | Count | Buyer Name | Delivery Type | Address | Date Time | Status | PayPal token | Comment |
|----|---------|-------|------------|---------------|---------|-----------|--------|--------------|---------|

It's filled out by bot, user shouldn't modify any data here.
User can find in this table some useful info about customer like his name, chosen delivery type, address or comment.
