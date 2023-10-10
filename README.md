# Dirry

> The Director Utility

This is a tool I've built to help extract resources from director files.

It's very work-in-progress and is probably a little broken.

It supports your usual director files, cast files, and self-contained EXE. (afterburner or not)

[![X Follow](https://img.shields.io/twitter/follow/_markeh)](https://twitter.com/_markeh)
[![Discord](https://img.shields.io/discord/368096211502366740?style=flat&logo=discord)](https://discord.gg/JTPQDFbhmv)

## Usage

### Dump

The dump utility can work on any director files, including embedded executable files.

```
dirry dump path/to/dir
```

This will dump everything into `~/dirry/out/dump`, to debug palette issues there is a `\_debug`` folder in there too with HTML output of them.

Note at the moment the application will verbosely log into the "logs" directory. If it gets too big just delete it.

### Zip

```
dirry zip path/to/dir
```

Will create a direcotr zip in out/zip

## Whats missing/broken?

### Some inputs

Mac resource fork files are planned but not done.
Windows NE exe's (PE is ok) might have issues too but I cannot get Director 4 to save an exe properly in Windows 98.

### Sounds

I just haven't gotten around to it yet.

### Bitmaps with 16bits

Also haven't gotten around to it, and my existing test files didn't have any. Expect it soon enough.

### Text XMED

Needs more research.

### Sometimes Cast Issues

I have gone through many projects about casts and I always seem to extract the common properties wrong of it wrong (slight offset issue I think).

### 3D Conversion

I do intend to do this but it looks like something done with Intel 3D’s IFX Toolkit

## Debugging

### Print Output

To see output on the screen you should use the `--verbose` flag on the dirry command, example:

```
dirry --verobse=all dump ..
```

Instead of "all" you can sub in the section you want to look at.

### Log Output

As the logs get quite big they are by default turned off. But if you would like to enable it to see whats happening,
or to create an issue, add the `--logging` flag.

```
dirry --logging dump ...
```

In `~/logs/` a new directory is created with the log file, and then each category/section split up apart (so you'll)
have double the logs, and again these can be super big so clean up!)

## New shockwave utilities

Currently I'm moving a lot of the afterburner stuff around, and have broken it. It won't work at the moment.
