# /bin/sh
mkdir dist/
echo "Building darwin"
GOOS=darwin GOARCH=amd64 go build
mv utm dist/utm_darwin_amd64

echo "Building linux"
#GOOS=linux GOARCH=386 go build
#mv aroz_online build/aroz_online_linux_i386
GOOS=linux GOARCH=amd64 go build
mv utm dist/utm_linux_amd64
GOOS=linux GOARCH=arm GOARM=6 go build
mv utm dist/utm_linux_arm
GOOS=linux GOARCH=arm GOARM=7 go build
mv utm dist/utm_linux_armv7
GOOS=linux GOARCH=arm64 go build
mv utm dist/utm_linux_arm64

echo "Building windows"
GOOS=windows GOARCH=amd64 go build
mv utm dist/utm_windows_amd64.exe

echo "OK"
