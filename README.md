Aphid
=====

Helpful hints for Catkin/CMake output

Installing
----------

Make sure you have Go installed. For example, on Ubuntu and Debian:

    sudo apt-get install golang-go

Then get aphid using the go tool:

	go get github.com/ijt/aphid

Running
-------

	cmake ../src 2>&1 | aphid

Aphid will insert helpful hints for errors it recognizes in the output of CMake.
These hints are marked with a prefix of [aphid].

To see an example, try
    
    aphid < ./test_data/could_not_find_a_config.txt

