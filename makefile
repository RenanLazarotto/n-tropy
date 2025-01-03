install: 
	@go build -o ntropy main.go
	@sudo cp ntropy /usr/local/bin