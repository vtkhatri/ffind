#include <stdlib.h>
#include <stdio.h>
#include <string.h>
#include <errno.h>

#define SHORT_ARG      0
#define FILENAME_ARG   1
#define PATH_ARG       2
#define CMD_ARGS       3

static FILE *logout;

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
	char *exec = (char *) malloc (sizeof(char));

	int isRegex = 0, isCaseInsensitive = 0, isMaxDepth = 0, isEqual = 0, isExpressionDone = 0, isPathDone = 0;

	for (int i=0; argv[i] != NULL; i++) {
		if (argv[i][0] == '-') {
			if (argv[i][1] == '-') {
				// it's a longArg
				if (strcmp(argv[i], "--debug") == 0) {
					fclose(logout);
					logout = stdout;
				} else if (strcmp(argv[i], "--help") == 0) {
					printUsage();
					exit(0);
				}
			} else {
				for (int j=1; j<strlen(argv[i]); j++) {
					switch (argv[i][j]) {
					case 'f':
						concatenateString(dashflags, " -type f");
						break;
					case 'd':
						concatenateString(dashflags, " -type d");
						break;
					case 'r':
						isRegex = 1;
						break;
					case 'i':
						isCaseInsensitive = 1;
						break;
					case 'e':
						isMaxDepth = 1;
						break;
					case '=':
						isEqual = 1;
						if (!isMaxDepth) {
							printf("ffind: =maxdepth used without -e flag\n");
							printUsage();
							exit(1);
						} else {
						}
						break;
					
					default:
						break;
					}
				}
			}
		} else {
			if (i != 0) {
				if (isExpressionDone) {
					if (isPathDone) {
						printf("ffind: too many arguments with no flags, only 2 accepted, <expression> followed by <path>\n");
						printUsage();
						exit(0);
					} else {
						isPathDone = 1;
						concatenateString(path, " ");
						concatenateString(path, argv[i]);
					}
				} else {
					isExpressionDone = 1;
					concatenateString(expression, " ");
					concatenateString(expression, argv[i]);
				}
			}
		}
	}

	if (isMaxDepth && !isEqual) {
		printf("ffind: -e flag used without =maxdepth\n");
		printUsage();
		exit(1);
	}

	if (isRegex) {
		if (isCaseInsensitive) concatenateString(dashflags, " -iregex");
		else concatenateString(dashflags, " -regex");
	} else {
		if (isCaseInsensitive) concatenateString(dashflags, " -iname");
		else concatenateString(dashflags, " -name");
	}

	concatenateString(command, path);
	concatenateString(command, dashflags);
	concatenateString(command, expression);
	concatenateString(command, exec);

	return 0;
}

int main(int argc, char **argv) {
	logout = fopen("/dev/null", "w");
	char *cmd = (char *) malloc (sizeof(char));
	int retval = makeCommand(cmd, argv);
	printf("cmd = %s\n", cmd);

	return retval;
}
