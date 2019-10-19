CC = go
TARGET = gobra
bindir = /usr/local/bin

$(TARGET): *.go
	$(CC) get -x github.com/rthornton128/goncurses
	$(CC) get -x github.com/hajimehoshi/oto
	$(CC) build -x
.PHONY all: $(TARGET)

.PHONY install: all
	mv $(TARGET) $(DESTDIR)$(bindir)/$(TARGET)

.PHONY uninstall:
	rm -f $(DESTDIR)$(bindir)/$(TARGET)

.PHONY clean:
	rm -f $(TARGET)
