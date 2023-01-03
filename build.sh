PWD=`pwd`
INSTALL=${PWD}"/bin"


echo "build linux ..."
cd ${PWD}
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags='-w -s' -o ${INSTALL}/linuxapp
echo "build mac ..."
cd ${PWD}
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags='-w -s' -o ${INSTALL}/macapp
