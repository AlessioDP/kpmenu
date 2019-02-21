PKGNAME := kpmenu
BINNAME := ${PKGNAME}
DESTDIR := /usr

.PHONY: all build clean install run uninstall

all: build

build:
	go get -v
	go build -v -o ${BINNAME}

clean:
	rm -f ${BINNAME}

install:
	install -Dm755 ${BINNAME} ${DESTDIR}/bin/${BINNAME}
	install -Dm644 LICENSE ${DESTDIR}/share/licenses/${PKGNAME}/LICENSE

run:
	go run main.go

uninstall:
	rm -f ${DESTDIR}/bin/${BINNAME}
	rm -f ${DESTDIR}/share/licenses/${PKGNAME}/LICENSE
