cd icon
cmd /c make_icon.bat iconwin.ico
cd ..
go install github.com/rakyll/statik
statik -f -src=static && go build -ldflags -H=windowsgui
