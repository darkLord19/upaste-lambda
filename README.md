# upaste-lambda

create pastes of data using aws lambda and store in dynamodb

# How to deploy
```
GOOS=linux CGO_ENABLED=0 go build main.go
zip main.zip main
```
Upload the zip to lambda function Code tab
