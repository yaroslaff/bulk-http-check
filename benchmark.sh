#!/bin/bash

FILELIST=/tmp/urls.txt

if [ -x ./bulk-http-check ]
then
    BULK=./bulk-http-check
else
    BULK=`which bulk-http-check`
fi

echo binary: $BULK

if [ ! -f $FILELIST ]
then
    echo missing $FILELIST
    exit 1
fi

for n in 1 10
do
    echo === $n connections ===
    LOG=/tmp/bulk-benchmark-$n.log
    RESULT=/tmp/bulk-benchmark-$n
    $BULK -n $n -b 10 -x 60 < /tmp/urls.txt > $LOG 2> $RESULT
    cat $RESULT
    echo ERRORS: `cat $LOG | grep " ERR " | wc -l`
    echo Statuses:
    grep " OK " < $LOG | cut -f 3 -d" "|sort|uniq -c|sort -n    
    echo
done
