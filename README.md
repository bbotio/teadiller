# teadiller


To run golang inside of docker container. You need to run the following command in project folder.
```
docker run --rm -it -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" golang bash
```
OR
compile and run
```
docker run --rm -t -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" --env TELEGRAMM_TOKEN="<your_token>" --env TELEGRAM_BOT_NAME="<your_botname>" --env BOT_WEB_PORT="<your_web_port>" -p $BOT_WEB_PORT:$BOT_WEB_PORT golang go run main.go
```
Run detached 
```
docker run -d -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" --env TELEGRAMM_TOKEN="<your_token>" --env TELEGRAM_BOT_NAME="<your_botname>" --env BOT_WEB_PORT="<web_port>" -p $BOT_WEB_PORT:$BOT_WEB_PORT golang go run main.go
```

## Configurations
| Env variable name | Value |
|-------------------|-------|
| TELEGRAMM_TOKEN   |\<Put telegram token here\>|
| TELEGRAMM_BOT_NAME|\<Put telegram bot name\>|
| BOT_WEB_PORT      |\<Put web port\>|

## Instagram api examples
```
msg.ReplyMarkup = tgbotapi.ReplyKeyboardMarkup{[][]string{{"black tea", "green tea", "t1", "t2", "t3", "t4", "t5", "t6", "t7", "t8"}}, false, false, false}
```
