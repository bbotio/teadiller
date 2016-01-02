# teadiller


To run golang inside of docker container. You need to run the following command in project folder.
```
docker run --rm -it -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" golang bash
```
OR
compile and run
```
docker run --rm -v "$PWD":/go/src/teadiller -w /go/src/teadiller --env GOPATH="/go" --env TELEGRAMM_TOKEN="<your_token>" golang go run main.go
```
