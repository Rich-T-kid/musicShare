outDirectory="protobuff"
client="client.go"
server="server.py"
generate:
	mkdir -p $(outDirectory)
#Generate Go client
	protoc --go_out=./$(outDirectory) --go-grpc_out=./$(outDirectory) --proto_path=. *.proto
# Generate Python server
	python3 -m grpc_tools.protoc --python_out=./$(outDirectory) --grpc_python_out=./$(outDirectory) --proto_path=. *.proto
# Implementation file for Go client
	touch $(outDirectory)/grpc_implementation.py
# Create client and server files
	touch $(client) && touch $(server)
# 
clean:
	rm -rf $(outDirectory) 