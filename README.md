# lcg

lastfm-chart-generator

**Requirements:**

Make sure you have golang correctly installed on your system. You may also want to `$GOPATH/bin` to your environment's `$PATH`.

**Install:**

Download the source files

`$ go get -u -v github.com/getmicah/`

Compile the binaries

`$ go install github.com/getmicah/`

**Usage:**

Assuming you added `$GOPATH/bin` to your `$PATH`, you should be able to run the program now.

`$ lcg <user> <period> <size>`

* user: <last.fm username>
* period: week, month, 3month, 6month, year, overall
* size: 3, 4, 5

**Example:**

    $ lcg getmicah week 3
    Loading user...
    Fetching covers...
    Drawing chart...
    Saved collage.png

![chart](collage.png?raw=true&s=300 "Chart")
