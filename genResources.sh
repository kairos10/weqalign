#!/bin/bash

outFile=temp_web.go
#files=( web/index.html web/main.js )
files=( web/*.html web/*.js web/*.jpg )

echo generating $outFile...
echo package main > $outFile

echo var _genFiles = map[string]string\{ >> $outFile
for f in ${files[@]}; do
	if [ -f $f ]; then
		echo adding $f...
		echo -n \"$f\": \` >> $outFile
		cat $f | base64 >> $outFile
		echo \`, >> $outFile
	fi
done 
echo \} >> $outFile
