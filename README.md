# sprm - A Simple File Renaming Program

This is a Go version of my spacerm program that is used to remove spaces 
and other undesirable characters from filenames. Files may either be renamed 
in place, or a copy of the file will be made with the new file name. Other 
options include verbose output that prints the old file name and the new 
file name including what operation was conducted. Also, the interactive 
option will prompt the user for confiration before copying or renaming a 
file. Extensions are never modified. 

**Getting Started**

First get a copy of this repository either with Git or using GitHub's download
option:

    git clone https://github.com/spcnvdr/sprm.git

Then change into the source code directory:

    cd sprm/src

Then build the program into an executable binary using Go:

    go build sprm.go

Finally, run the program to get a list of options:

    ./sprm --help


**Examples**

Remove spaces from 'file 1' changing it to 'file1'. Repeated for the
other files specified.

    sprm 'file 1' 'file 2' 'file 3'

Rename '/home/user/bad file name.pdf' to '/home/user/bad-file-name.pdf'
and prompt the user to confirm this before renaming.

    sprm -d -i /home/user/bad\ file\ name.pdf

Rename the file 'bad.file.name.pdf' to 'badfilename.pdf'

    sprm --strip=. bad.file.name.pdf


**To-Do**

- [ ] - Use os.Stat to determine if input file is a normal file.
- [ ] - Do nothing if output file already exists

**Contributing**

Pull requests, new feature suggestions, and bug reports/issues are
welcome.


**License**

This project is licensed under the 3-Clause BSD License also known as the
*"New BSD License"* or the *"Modified BSD License"*. A copy of the license
can be found in the LICENSE file. A copy can also be found at the
[Open Source Institute](https://opensource.org/licenses/BSD-3-Clause)

