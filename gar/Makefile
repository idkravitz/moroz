.PHONY: clean target

#include $(GOROOT)/src/Make.inc
GOOS=windows
O=8

ifeq ($(GOOS),windows)
clean=del /Q
else
clean=rm -rf
endif


target:
	gobuild -o bin/gar.exe

clean:
	$(clean) *.$(O) bin