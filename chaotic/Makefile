all: 	clean bindata.go

clean:
	rm -f bindata.go

bindata.go:
	go-bindata -pkg=chaotic -prefix=public public
