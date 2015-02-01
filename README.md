# MW
(as in "MiddleWare")

[![Build Status](https://travis-ci.org/morhekil/mw.svg?branch=master)](https://travis-ci.org/morhekil/mw)

This is a collection of various Go middleware that I use in my projects.
For non-trivial middlewares, the detailed READMEs are available in their
directories. Feel free to open an issue if you need help, or found a bug
in any of those.

## [Chaotic](https://github.com/morhekil/mw/tree/master/chaotic)

Chaotic provides stdlib-compatible middleware to inject configurable
delays and failures into the requests processed by its underlying HTTP stack.
