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

int makeCommand(char *command, char** argv) {

	for (int i=0; argv[i] != NULL; i++) {
		if (argv[i][0] == '-') {
			if (argv[i][1] == '-') {
				// it's a longArg
				if (strcmp(argv[i], "--debug") == 0) {
					strncat(command, "debug ", 6);
				} else if (strcmp(argv[i], "--help") == 0) {
					strncat(command, "help ", 5);
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
	return 0;
}

int main(int argc, char **argv) {

	char *cmd = (char *) malloc (sizeof(char)*20);
	int retval = makeCommand(cmd, argv);
	printf("cmd = %s\n", cmd);
	return retval;
}
