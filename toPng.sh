#!/bin/bash

# This script convert images to png format

# Verify that mogrify is available

set -e

if [ -z "$(which mogrify 2> /dev/null)" ]; then
	echo 'Error: This script requires mogrify. Install it with `apt-get install imagemagick`'
	exit 1
fi

mogrify -format png $1
