#!/bin/bash

set -e

mkdir -p images/input images/output images/pattern

docker run --rm -v "$(pwd)/images":"/images" twikle $@
