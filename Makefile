PREFIX?=/usr/local
MANPREFIX=${PREFIX}/share/man

all: slideshow slideshow.1

slideshow:
	go build entf.net/slideshow

slideshow.1: slideshow.1.scd
	scdoc < $< > $@

install: all
	mkdir -p ${DESTDIR}${PREFIX}/bin
	cp -f slideshow ${DESTDIR}${PREFIX}/bin
	mkdir -p ${DESTDIR}${MANPREFIX}/man1
	cp -f slideshow.1 ${DESTDIR}${MANPREFIX}/man1

uninstall:
	-rm -f ${DESTDIR}${PREFIX}/bin/slideshow
	-rm -f ${DESTDIR}${MANPREFIX}/man1/slideshow.1

.PHONY: all slideshow install uninstall
