# go-browse

An attempt at an HTML/CSS rendering engine, written in Go.

## Background

This is heavily based on Matt Brubeck's excellent "Let's build a browser engine!" article series published on [his website](https://limpet.net/mbrubeck/2014/08/08/toy-layout-engine-1.html) and was built mostly as an educational exercise during my time at [Recurse Center](http://recurse.com/) in early 2020.

This implementation is currently (and may remain) incomplete. Here are some features that I would like to add in the future.

- [ ] Painting the LayoutTree onto a canvas
- [ ] Supporting varying colors of elements
- [ ] Support for Inline CSS Box elements, including spans and text
- [ ] Support for percentage-based length calculations
- [ ] Support for partial HTML
- [ ] Correct error-handling for incorrect CSS rules/properties

## Instructions

From the root folder, first install all dependencies.

```bash
  go install ./...
```

Then use the `run` command in the Makefile.

```bash
  make run
```
