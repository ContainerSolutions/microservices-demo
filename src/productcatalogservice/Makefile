CUSTOM_GOPATH=$(shell pwd)/../..

productcatalogservice:
	GOPATH=$(CUSTOM_GOPATH) go build -o $@ server.go database.go

.PHONY:run
run: productcatalogservice 
	DB_USER=root DB_PASS=mypassword DB_HOST=127.0.0.1 DB_PORT=3306 DB_NAME=mydb ./productcatalogservice

.PHONY:mysql
mysql: 
	docker run -P --rm -d --net=host --name=mysqltest \
		--env="MYSQL_ROOT_PASSWORD=mypassword" \
		--env="MYSQL_DATABASE=mydb" \
		mysql || true

.PHONY:dep
dep:
	export GOPATH=$(CUSTOM_GOPATH)
	dep ensure --vendor-only

clean:
	rm -f productcatalogservice
	killall productcatalogservice || true
	docker kill mysqltest || true
