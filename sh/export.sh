#!/bin/bash
height=$1
node=$2
./bacd export      --height  ${height}  >  bcv_${node}_${height}.json
