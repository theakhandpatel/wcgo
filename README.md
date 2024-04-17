# Word Count in Go (WCGO)
This utility is a Go implementation of the Unix wc command. It provides statistics about a file or input from the standard input (stdin). The statistics include the number of bytes, lines, words, and characters.

## Usage
The utility can be run from the command line with various flags to specify what statistics to output. If no flags are provided, it defaults to counting bytes, lines, and words.

Here are the available flags:

* -c: Count the number of bytes in the input.
* -l: Count the number of lines in the input.
* -w: Count the number of words in the input.
* -m: Count the number of characters in the input.


You can specify one or more flags. For example, to count the number of lines and words in a file, you would use the -l and -w flags:

```bash
 wcgo -l -w myfile.txt
  ```


### Standard Input
If no file is specified, the utility reads from stdin. For example, you can pipe the output of another command into this utility:

```bash
 echo "Hello, world!" | wcgo -l -w 
```

### Multiple Files
This utility also supports multiple files. If more than one file is specified, it will output statistics for each file and a total count:

```bash
 wcgo test.txt main.go 
```

Output
The output is a string that includes the counts specified by the flags, followed by the file name. If an error occurs while processing a file, the output will include the error message.

For example, the output might look like this:

```
  7145  58164 342190 test.txt
   155    413   3569 main.go
  7300  58577 345759 total
```

The first number is the number of lines, the second number is the number of words, and the third number is the number of bytes. The file name is the last item in the output.

## Usage
 Pre-requisites: Go should be installed on the system. If not, you can download it from [here](https://golang.org/dl/). Make tools should also be installed on the system.

Follow these steps:

1. **Clone the repository:**
```bash
git clone https://github.com/theakhandpatel/wcgo.git
```

2. Navigate to the project directory:
```bash
cd wcgo
```



### Run (Without Installation)
To run the utility, use the following command:
```
go run main.go -l -w test.txt
```

### Install


1. Install the utility using makefile:
```bash
make install
```

2. Run the utility with the desired flags:
```bash
wcgo -l -w test.txt
```

### Uninstall Instructions
To uninstall the utility, run the following command:
```bash
make uninstall
```


## Inspiration
[Coding Challenge: Word Count](https://codingchallenges.fyi/challenges/challenge-wc/)

