PKGNAME := kpmenu
BINNAME := ${PKGNAME}

.PHONY: all build clean install run uninstall

all: build

build:
	go build -v -o ${BINNAME}

clean:
	rm -f ${BINNAME}

install:
	install -Dm755 ${BINNAME} /usr/bin/${BINNAME}
	install -Dm644 LICENSE /usr/share/licenses/${PKGNAME}/LICENSE

run:
	go run main.go

uninstall:
	rm -f /usr/bin/${BINNAME}
	rm -f /usr/share/licenses/${PKGNAME}/LICENSE
