# osutil

Access to operating system functionality dependent of every platform and
utility packages for the Shell.

+ config/env: set persistent environment variables
+ config/shconf: parser and scanner for the configuration in format shell-variable
+ distro: detects the Linux distribution
+ file: common operations in files
+ pkgutil: basic operations for the management of packages in operating systems
+ sh: interprets a command line like it is done in the Bash shell
+ user: provides access to UNIX users database in local files
+ user/crypt: password hashing used in UNIX

[Documentation online](http://godoc.org/github.com/tredoe/osutil)

## Testing

`go test ./...`

`sudo env PATH=$PATH go test -v ./...`

'sudo' command is necessary to copy the files '/etc/{passwd,group,shadow,gshadow}' to the temporary directory, where the tests are run.
Also, it uses 'sudo' to check the package manager, at installing and removing the package 'mtr-tiny'.


## License

The source files are distributed under the [Mozilla Public License, version 2.0](http://mozilla.org/MPL/2.0/),
unless otherwise noted.  
Please read the [FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html)
if you have further questions regarding the license.
