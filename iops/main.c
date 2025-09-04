#include <stdio.h>
#include <stdlib.h>

#define BUF_SIZE 4096   // 4KB page-sized buffer

int main(int argc, char *argv[]) {
    if (argc < 2) {
        fprintf(stderr, "Usage: %s <filename>\n", argv[0]);
        return 1;
    }

    FILE *fp = fopen(argv[1], "rb");
    if (!fp) {
        perror("fopen");
        return 1;
    }

    char *buf = malloc(BUF_SIZE);
    if (!buf) {
        perror("malloc");
        fclose(fp);
        return 1;
    }

    size_t n;
    while ((n = fread(buf, 1, BUF_SIZE, fp)) > 0) {
        // Do nothing, just read
    }

    free(buf);
    fclose(fp);

    return 0;
}
