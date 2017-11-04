# lastfm-chart-generator

Dependencies:

`$ brew install pkg-config cairo pango libpng jpeg giflib`

Install:

`$ npm install -g getmicah/lastfm-chart-generator`

Usage:

`$ lcg <user> <period> <size>`

* **user**: last.fm username
* **period**: week, month, 1month, 3month, 6month, year, overall
* **size**: 3, 4, 5

Example:

    $ lcg getmicah week 3
    created collage.png