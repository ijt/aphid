Aphid
=====

Helpful hints for Catkin/CMake output

Installing
----------

	go get github.com/ijt/aphid

Running
-------

	cmake ../src 2>&1 | aphid

Aphid will insert helpful hints for errors it recognizes in the output of CMake.
These hints are marked with a prefix of [aphid].

