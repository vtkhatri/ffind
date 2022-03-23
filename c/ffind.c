#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

#define SHORT_ARG      0
#define FILENAME_ARG   1
#define PATH_ARG       2
#define CMD_ARGS       3

void printUsage() {
	printf("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]\n");
	return;
}

int concatenateString(char *dest, char *src) {
	if (strlen(src) == 0) return 0;

	if (dest == NULL) {
		dest = (char *) malloc (strlen(src)*sizeof(char));
	} else {
		dest = realloc (dest, strlen(dest) + strlen(src));
	}
	strncat(dest, src, strlen(src));
	return 0;
}

int makeCommand(char *command, char** argv) {
	char *expression = (char *) malloc (sizeof(char));
	char *path = (char *) malloc (sizeof(char));
	char *dashflags = (char *) malloc (sizeof(char));
	char *doubledashflags = (char *) malloc (sizeof(char));

	for (int i=0; argv[i] != NULL; i++) {
		if (argv[i][0] == '-') {
			if (argv[i][1] == '-') {
				// it's a longArg
				if (strcmp(argv[i], "--debug") == 0) {
					concatenateString(doubledashflags, "debug ");
				} else if (strcmp(argv[i], "--help") == 0) {
					concatenateString(doubledashflags, "help ");
					printUsage();
					//exit(0);
				}
			} else {
				// it's a shortArg
				// cmd[SHORT_ARG].append(argv[i]);
			}
		} else {
			if (i != 0) {
				// strncat(command, "1", 1);
			}
		}
	}

	concatenateString(command, expression);
	concatenateString(command, path);
	concatenateString(command, dashflags);
	concatenateString(command, doubledashflags);

	return 0;
}

int main(int argc, char **argv) {
	char *cmd = (char *) malloc (sizeof(char));
	int retval = makeCommand(cmd, argv);
	printf("cmd = %s\n", cmd);

	return retval;
}
