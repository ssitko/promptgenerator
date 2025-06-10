#!/bin/bash

go run main.go \
    -a "Experienced golang developer" \
    -p '''
    
    Based on input program, generate documentation guide how to use it.

    ''' \
    -i ./main.go \
    -o ./README.md \
    -f ai.db
