build:
	go build -o bin/ovh-dynip main.go

clean:
	rm -f bin/ovh-dynip

all: clean build

install:
	sudo cp bin/ovh-dynip /usr/bin/
