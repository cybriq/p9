#!/bin/bash

go install ./pod/buidl/.
buidl generate
buidl install

buidl help
echo "buidl is the build system for p9. `buidl help` to see what it does"