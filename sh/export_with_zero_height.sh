#!/bin/bash
height=$1
node=$2
./bacd export  --for-zero-height    --height  ${height}  >  bcvwithzeroheight_${node}_${height}.json
