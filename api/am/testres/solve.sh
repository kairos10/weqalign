#!/bin/bash

keepSolvedDir=~/_solved

if [ $# -lt 1 ]; then
	echo Use $0 FILE.JPG [other solve-field parameters]
	exit 1
fi

fpName=$1
fileName=${fpName##*/}
if [ ${#fileName} != ${#fpName} ]; then
	pathName=${fpName%/*}
else
	pathName="."
fi
namePart=${fileName%.*}

shift

do_exit() {
        rm -f "${pathName}/${namePart}.working"

	if [ -f ${pathName}/${namePart}.solved ]; then
		exit 0
	else
		rm -f "${pathName}/${namePart}.working"
		touch "${pathName}/${namePart}.notsolved"
		exit 3
	fi
}
trap do_exit EXIT

# save the image, since it might get deleted through the cleanup process
if [ -n "$keepSolvedDir" ]; then
	if [ ! -d $keepSolvedDir ]; then mkdir -p $keepSolvedDir; fi
	if [ -f $fpName ]; then cp -n $fpName $keepSolvedDir/; fi
fi

if [ -f $fpName ]; then

	touch "${pathName}/${namePart}.working"
	rm -f "${pathName}/${namePart}.notsolved"

	set -x
	(time solve-field --overwrite --downsample 1 --depth 30 --cpulimit 90 --resort --no-remove-lines --uniformize 0 --objs 50 --crpix-center --no-plots --new-fits none --rdls none --corr none --match none --index-xyls none --dir $pathName $fpName $*) 2>&1 | tee -a ${pathName}/${namePart}.log
	set +x

	if [ -f ${pathName}/${namePart}.solved -a -n "$keepSolvedDir" ]; then
		if [ ! -d $keepSolvedDir ]; then mkdir -p $keepSolvedDir; fi
		cp -n ${pathName}/${namePart}.* $keepSolvedDir
	fi

fi
