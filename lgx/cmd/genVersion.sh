#!/bin/bash

# ----------------------------------------------------------------------------------
# genVersion.sh
# Copyright 2019,2020 by ETOS GmbH
# Author: Waldemar Urbas
#-----------------------------------------------------------------------------------
# History
# ----------------------------------------------------------------------------------
# 2019.12.23 (wu) Init
# ----------------------------------------------------------------------------------


vfile=/tmp/xVersion_$1.txt
now=$(date +'%y.%m.%d')

if [ -e $vfile ]; then
  v=`grep $now $vfile`
else
  v=$now.0  
fi

av=${v:9:2}
nVersion=$((av+1))
echo -n $now.$nVersion > $vfile

ja=$(date +'%Y')
jj=$((ja-2010))

nVersion="$jj.$(date +'%m.%d').$nVersion"
echo -n $nVersion 
