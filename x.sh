#!/bin/bash
#-------------------------------------------------------------
# 2020.07.01 (wu) Github-Version setzen
#-------------------------------------------------------------
echo
echo "Git-Version-Vergabe (c) 2020 by Wald.Urbas"
echo
git tag --sort=-committerdate | head -1 >/dev/null 2>&1
if [ $? -ne 0 ]; then
  echo "kein gitRepo"
  exit 1
fi

over=$(git tag --sort=-committerdate | head -1)
echo "oldVersion: $over"

ver=$1
if [ "$ver" == "" ]; then
  echo "neue version fehlt: z.B.: 1.5.1"
  exit 1
fi

ver=${ver[@]/v/}

echo "newVersion: v$ver"


over=${over[@]/v/}
over=${over[@]/./}
over=${over[@]/./}

nver=${ver[@]/./}
nver=${nver[@]/./}

#if [ $nver -gt $over ]; then
if [ $nver -eq $over ]; then
  echo "version update"

cdir=$(pwd)
mdir=/home/master/go/src
gdir=https://proxy.golang.org

url=${cdir[@]/$mdir/$gdir}/@v/v$ver.info
echo "$url"

echo 0: ${cdir}
echo 1: ${mdir}
echo 2: ${cdir[@]}
echo 3: ${cdir[@]/$mdir}
echo 4: ${cdir[@]/$mdir/$gdir}


#git tag v$ver
#git push origin v$ver


#echo "curl https://proxy.golang.org/github.com/waldurbas/got/@v/v1.6.1.info"

#curl $url
else
echo "no version update"
fi

echo
