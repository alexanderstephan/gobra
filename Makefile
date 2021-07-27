CC = go
TARGET = gobra
bindir = /usr/local/bin

$(TARGET): *.go
	# Fix for ncurses install: https://github.com/rthornton128/goncurses/issues/56
	export CGO_CFLAGS_ALLOW=".*"
	export CGO_LDFLAGS_ALLOW=".*"
	$(CC) get -x github.com/alexanderstephan/goncurses 
	$(CC) get -x github.com/hajimehoshi/oto
	$(CC) build -x
	export CGO_CFLAGS_ALLOW=
	export CGO_LDFLAGS_ALLOW=

all: $(TARGET)

install: all
	mv $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

uninstall:
	rm -f $(DESTDIR)$(bindir)/$(TARGET)

clean:
	rm -f $(TARGET)
