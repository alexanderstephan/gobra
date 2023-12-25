CC = go
TARGET = gobra
bindir = /usr/local/bin
SOURCE_DIRS := $(shell find . -type d)
vpath %.go $(SOURCE_DIRS)

$(TARGET): cmd/main.go
	# Fix for ncurses install: https://github.com/rthornton128/goncurses/issues/56
	export CGO_CFLAGS_ALLOW=".*"
	export CGO_LDFLAGS_ALLOW=".*"
	$(CC) get -x github.com/rthornton128/goncurses
	$(CC) get -x github.com/hajimehoshi/oto
	$(CC) build -o $(TARGET) -x $<
	export CGO_CFLAGS_ALLOW=
	export CGO_LDFLAGS_ALLOW=

.PHONY: all
all: $(TARGET)

.PHONY: install
install: all
	mv $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

.PHONY: uninstall
uninstall:
	rm -f $(DESTDIR)$(bindir)/$(TARGET)

.PHONY: clean
clean:
	rm -f $(TARGET)
