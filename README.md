# Podcast Studio

This project is being written as part of the SUSE Hackweek 10 [1].

I've been running a podcast together with two friends for over a year now. We have a recording and mixing process going, but unfortunately it involves using non-free operating systems and software all the way through.

It is probably going to be hard to get away from using the things we use for the actual recording part for a while, but one thing I feel should be perfectly doable in openSUSE is the mastering bit. In theory, Audacity could work for this, but it is old and clunky and doesn't do real-time filtering, making it annoying to use.

I also want a project to learn google Go.

My project is to write a new program, ideally similar to audacity in ultimate utility, but very limited in scope to begin with:

* Load individual speakers as separate tracks

* Load extra sounds like intro/outro jingles, to be layered over the speaking tracks

* Apply limiting / compression / equalizing to the tracks to level

* Apply the same effects to the master track

* Effects are chained and applied real-time, not as in audacity.

## TODO

* Investigate possible UI options, preferrably not tying the program to any specific desktop environment.

* Learn basic sound programming :]

* Hack!

  [1]: https://hackweek.suse.com/projects/104 "SUSE Hackweek 10"
