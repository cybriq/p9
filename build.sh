#!/bin/bash

go install ./pod/buidl/.
buidl generate
buidl install

buidl
echo "buidl is the build system for p9. run 'buidl' by itself to see what it does"
echo "currently shader recompilations don't work, but the bundled code works already"