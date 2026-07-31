build:
	go build -o lyricist
	@echo "Build passed"

clean:
	go clean
	rm lyricist
