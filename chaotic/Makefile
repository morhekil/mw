all: 	clean bindata.go

clean:
	rm -f ./bindata/bindata.go

bindata.go:
	go-bindata -pkg=bindata -o ./bindata/bindata.go -prefix=public public
