default: run

run: build
	./server

build: 
	go build -o server main.go

clean:
	rm server && rm log.txt

clear:
	clear

tables:
	cd dev-scripts/db && ./add_tables.sh && cd ../../