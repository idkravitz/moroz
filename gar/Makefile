# Tools set depends on arch
ifeq ($(shell arch), x86_64)

A=6

else

A=8

endif

GC=$(A)g
GL=$(A)l

all: gar

gar: main.6
	$(GL) -o $@ $<

%.$(A): %.go
	$(GC) $<
