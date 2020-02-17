#!/bin/bash

outFile=temp_web.go
files=( web/index.html web/main.js )

echo generating $outFile...
echo package main > $outFile

echo var _genFiles = map[string]string\{ >> $outFile
for f in ${files[@]}; do
	if [ -f $f ]; then
		echo adding $f...
		echo -n \"$f\": \` >> $outFile
		cat $f | sed -e 's?//.*??' >> $outFile
		echo \`, >> $outFile
	fi
done 
echo \} >> $outFile
