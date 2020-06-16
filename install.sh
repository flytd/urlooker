if [ "x$GOPATH" = "x" ];then
	echo "GOPATH is not set, please set $GOPATH!"
	exit 1
fi

export GO111MODULE="off"
echo "set GO111MODULE=off"
mkdir -p $GOPATH/src/github.com/flytd
   
if [ ! -d "$GOPATH/src/github.com/flytd/urlooker" ]; then
	cd $GOPATH/src/github.com/flytd && git clone https://github.com/flytd/urlooker.git
else
	cd $GOPATH/src/github.com/flytd/urlooker && git pull
fi
cd $GOPATH/src/github.com/flytd/urlooker && ./control.sh build
echo "install ok! run cd $GOPATH/src/github.com/flytd/urlooker && ./control.sh start all"