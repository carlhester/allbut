# allbut

Delete all files in your current directory except the files you identify.  

## Usage

$ allbut file-to-protect

- Given a current path with 4 files: testfile1, testfile2, testfile3, testfile4
- Run _allbut testfile1_ to remove all of the files and leave only testfile1 in current path
- Also works with multiple files _allbut testfile1 testfile2_
- Normal operation just prints the work to do, the -f flag is required to actually delete

