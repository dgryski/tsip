

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "tsip.h"

int main(int argc, char *argv[]) {

    const unsigned char seed[] ={ 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f };

    if (argc != 2) {
        printf("usage: %s testdata/tsip.txt\n", argv[0]);
        exit(1);
    }

    FILE *f = fopen(argv[1], "r");
    if (!f) {
        perror(argv[1]);
        exit(1);
    }

    size_t i = 0;
    unsigned char buf[64];
    uint64_t want;
    while (!feof(f) && fscanf(f, "%016llx", &want) == 1) {
        uint64_t got = tsip(seed, (const unsigned char *)buf, i);
        if (got != want) {
            printf("%lu: got=%016llx want=%016llx\n", i, got, want);
            exit(1);
        }
        buf[i]=i;
        i++;
    }

    fclose(f);

    printf("PASS\n");

    return 0;
}
