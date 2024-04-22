# Try the example app

## example1
```bash
cd example1
go run main.go
``` 

## example2
```bash
cd example2
go install github.com/rakyll/statik
statik -f -src=static -p statik
go run main.go
``` 

## example3
```bash
cd example3
go install github.com/rakyll/statik
statik -f -src=static -p statik
go run icon.go main.go
``` 

## example4
```bash
cd example4
go install github.com/rakyll/statik
statik -f -src=static -p statik
go run icon1.go icon2.go main.go
``` 

## example5
```bash
cd example5
go install github.com/rakyll/statik
statik -f -src=static -p statik
go run main.go
``` 
