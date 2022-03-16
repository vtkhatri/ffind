#include <stdlib.h>
#include <errno.h>

#define SHORT_ARG      0
#define FILENAME_ARG   1
#define PATH_ARG       2
#define CMD_ARGS       3

void printUsage() {
	printf("Usage: ffind [-fdri] [-e=maxdepth] [--debug --help] [expression] [path]");
	return;
}

char **makeCommand(char** argv) {

	char *cmd[CMD_ARGS];

	for (int i=0; argv[i] != NULL; i++) {
		if (argv[0] == '-') {
			if (argv[1] == '-') {
				// it's a longArg
				if (argv[i] == "--debug") {

					break;
				} else if (argv[i] == "--help") {
					printUsage();
					return NULL;
				}
			} else {
				// it's a shortArg
				// cmd[SHORT_ARG].append(argv[i]);
			}
		} else {
			if (i != 0) {

			}
		}
	}
	return cmd;
}

int main(int argc, char **argv) {

	if (makeCommand(argv) == NULL) {
		return 1;
	}
	return 0;
}
