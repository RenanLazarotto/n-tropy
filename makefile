install: 
	@go build -o ntropy main.go
	@cp ntropy ~/.local/bin
	