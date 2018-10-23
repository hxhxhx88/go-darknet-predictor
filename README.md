Here I provide a GoLang wrapper for my [forked Darknet](https://github.com/hxhxhx88/darknet), in which some [forward APIs](https://github.com/hxhxhx88/darknet/blob/master/src/forward.h) are provided.

To compile the code, first compile the above Darknet, then make sure your `libdarknet.so` can be linked through `-ldarknet`. One way to achieve this is to soft-link it to `/usr/local/lib`.
