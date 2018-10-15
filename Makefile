test:
	-dep init
	dep ensure
	go test ./controller
	go test ./allocator

clean:
	rm -rf ./vendor
