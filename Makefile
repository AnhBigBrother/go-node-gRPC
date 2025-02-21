gen_proto:
	cd server; sh gen_proto.sh;
	cd ..;
	cd client; sh gen_proto.sh;
	cd ..;

cert: 
	cd cert; sh gen.sh; cd ..

setup: 
	cd server; sh gen_proto.sh; go mod tidy;
	cd ..;
	cd client; sh gen_proto.sh; npm install;
	cd ..;
	cd cert; sh gen.sh; 
	cd ..;

server: 
	cd server; go run main.go
client: 
	cd client; ts-node main.ts
stream: 
	cd client; ts-node main.ts --cmd=bidirection_stream


.PHONY: gen_proto cert setup server client stream