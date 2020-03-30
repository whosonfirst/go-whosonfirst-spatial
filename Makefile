CWD=$(shell pwd)

go-bindata:
	mkdir -p cmd/go-bindata
	mkdir -p cmd/go-bindata-assetfs
	curl -s -o cmd/go-bindata/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata/master/cmd/go-bindata/main.go
	curl -s -o cmd/go-bindata-assetfs/main.go https://raw.githubusercontent.com/whosonfirst/go-bindata-assetfs/master/cmd/go-bindata-assetfs/main.go

bake: bake-static bake-templates

bake-static:
	go build -mod vendor -o bin/go-bindata cmd/go-bindata/main.go
	go build -mod vendor -o bin/go-bindata-assetfs cmd/go-bindata-assetfs/main.go
	rm -f static/*~ static/css/*~ static/javascript/*~
	@PATH=$(PATH):$(CWD)/bin bin/go-bindata-assetfs -pkg http static/javascript static/css

bake-templates:
	mv bindata.go http/assetfs.go
	rm -rf templates/html/*~
	bin/go-bindata -pkg templates -o assets/templates/html.go templates/html

debug:
	@make bake
	go run -mod vendor cmd/server/main.go -nextzen-apikey $(APIKEY) -api
