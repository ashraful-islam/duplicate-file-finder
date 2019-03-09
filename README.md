# duplicate-file-finder
a small prototype golang program to find and report duplicate files

## Disclaimer

This is a learning project, hence it is not recommended for usage in important/serious environments.
The author of this project do not hold any responsibility/liabilities for anykind of damage or loss resulted from using this project.
Use at your own risk. 

## How To:

Written in: Golang `v1.12`

First install Golang, instructions can be found [here](https://golang.org/doc/install) from Golang docs.

Then clone this repository 

In console/terminal, use following command to clone the repository
```$bash
git clone https://github.com/ashraful-islam/duplicate-file-finder.git
```

To build

```$bash
go build .
```

To run tests
```$bash
go test .
```

## Usage

The program(once built), will scan through a given root directory
into all sub-directories and check all files. Then it will check for duplicate files.

To run a scan over the directory "/home/me/downloads" and find all duplicate files in `download` directory,
use the following command in terminal/console

```bash
./duplicate-file-finder -path /home/me/downloads
```
For windows, use the .exe file as following(given this directory is in C:\ drive)
```
duplicate-file-finder.exe -path c:\home\me\downloads
```

A shorthand version of flag `-path` is `-p`, for example:
```
./duplicate-file-finder -p /home/me/downloads
```

Once the program has completed, it will report a list of duplicate files and where they are,
and a mini statistics of number of files and sizes(in bytes)

For example, using the test_data directory run gives following result:
```
[info] Searching for duplicates in:  test_data/
[status] Processing (this may take some time)
[status] Checking for duplicates (this will take some time)
----
Duplicate Group:  2
Name:  sample.txt
Size:  36
Path:  test_data/sample.txt
Hash:  d9e944f9126aaf2c9f68c10e907800da

Name:  sample_dup.txt
Size:  36
Path:  test_data/sub_dir/sample_dup.txt
Hash:  d9e944f9126aaf2c9f68c10e907800da

[status] Done!
Result Statistics:
-- Scanned:
Files Found: 3
Size(bytes): 112
-- Duplicates:
Files Found: 2
Size(bytes): 72
```

Here, we see that two duplicate files `sample.txt` and `sample_dup.txt` were found.
The program found total `3` files of cumulative `112 bytes` of data. 
It found, `2` duplicate files of cumulative `72 bytes` of data.
The `Duplicate Group: 2` indicates that there were two exactly same files in one group. 

## Todo:

- [x] add basic tests
- [ ] add more comments
- [ ] refactor
- [ ] add more tests