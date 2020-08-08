CFLAGS=-O2

.c.o:
	gcc -c $(CFLAGS) $< -o $@

tsiptest: tsip_test.o tsip.o
	gcc $(CFLAGS) $^ -o $@

test: tsiptest
	./tsiptest ../go/testdata/tsip.txt
