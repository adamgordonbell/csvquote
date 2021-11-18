# csvquote

_smart and simple CSV processing on the command line_

Are you looking for a way to process CSV data with standard UNIX shell commands?

Are you running into problems with embedded commas and newlines that mess
everything up?

Do you wish there was some way to add some CSV intelligence to these UNIX tools?

* awk, sed
* cut, join
* head, tail
* sort, uniq
* wc, split

This program can be used at the start and end of a text processing pipeline
so that regular unix command line tools can properly handle CSV data that
contain commas and newlines inside quoted data fields.

Without this program, embedded special characters would be incorrectly
interpretated as separators when they are inside quoted data fields.

By using csvquote, you temporarily replace the special characters inside quoted
fields with harmless nonprinting characters that can be processed as data by
regular text tools. At the end of processing the text, these nonprinting
characters are restored to their previous values.

In short, csvquote wraps the pipeline of UNIX commands to let them work on
clean data that is consistently separated, with no ambiguous special
characters present inside the data fields.

By default, the program expects to use these as special characters:

    " quote character  
    , field delimiter  
    \n record separator  

It is possible to specify different characters for the field and record
separators, such as tabs or pipe symbols.

Note that the quote character can be contained inside a quoted field
by repeating it twice, eg.

    field1,"field2, has a comma in it","field 3 has a ""Quoted String"" in it"

Typical usage of csvquote is as part of a command line pipe, to permit
the regular unix text-manipulating commands to avoid misinterpreting
special characters found inside fields. eg.

    csvquote foobar.csv | cut -d ',' -f 5 | sort | uniq -c | csvquote -u

or taking input from stdin,

    cat foobar.csv | csvquote | cut -d ',' -f 7,4,2 | csvquote -u

other examples:

    csvquote -t foobar.tsv | wc -l

    csvquote -q "'" foobar.csv | sort -t, -k3 | csvquote -u

    csvquote foobar.csv | awk -F, '{sum+=$3} END {print sum}'

## Installation

To install using go, clone the repo and then build it:
```
> go build -o csvquote cmd/cvsquote/main.go
> cp ./csvquote /usr/local/bin
```

To install using Earthly (linux):
```
earthly github.com/adamgordonbell/csvquote+build
cp ./csvquote /usr/local/bin
```

Install on MacOS (X86)
```
earthly github.com/adamgordonbell/csvquote+for-darwin-amd64
cp ./csvquote /usr/local/bin
```

Install on Windows
```
earthly github.com/adamgordonbell/csvquote+for-windows-amd64
# Then add to your path
```

## History

This is a fork of original version in C by Dan Brown found here: 
https://github.com/dbro/csvquote

More specifically, this is a fork of a fork of the C version. It is based on a GoLang version Dan wrote at some point. Go makes cross-compilation easier.
