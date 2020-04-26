.PHONY: proto manager image clean

manager:
	go build -o build/manager
	
proto:
	protoc -I proto proto/container/container.proto --go_out=plugins=grpc:pb/container
	protoc -I proto proto/worker/worker.proto --go_out=plugins=grpc:pb/worker

image:
	docker build -t jointfaas/manager .

clean:
	rm -rf build/*