catkin sleuth
=============

Helpful hints for Catkin/CMake output

Installing
----------

	go get .

or

	go get github.com/ijt/catkin_sleuth

Running
-------

	cmake ../src 2>&1 | catkin_sleuth

Catkin sleuth will insert helpful hints for errors it recognizes in the output of CMake.
These hints are marked with a prefix of [catkin_sleuth].

