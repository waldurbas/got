ver=$(git tag --sort=-committerdate | head -1)
echo "Version: $ver"

ver=${ver[@]/v/}

# in $cdir $sdir mit waldurbas ersetzen
cdir=$(pwd)
sdir=$GOPATH/src/wux
ndir=${cdir/$sdir/waldurbas}

# in $ndir $sdir mit waldurbas ersetzen
sdir=$GOPATH/src/etos
ndir=${ndir/$sdir/etos-group}


url=https://proxy.golang.org/$ndir/@v/v$ver.info
echo
echo "URL: $url"
