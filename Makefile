APPNAME		:= citb

build:
	go build -ldflags '-X main.token=$(TOKEN)'

clean:
	rm -f $(APPNAME)
