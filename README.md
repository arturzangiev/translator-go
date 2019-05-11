Basically, the program takes every entry from source.csv translate it and generate new file output.csv with translations.

# To run the program
go run -ldflags "-X main.apiKey=INSERT_API_KEY_HERE" translator.go

# To build the program
go build -ldflags "-X main.apiKey=INSERT_API_KEY_HERE" translator.go