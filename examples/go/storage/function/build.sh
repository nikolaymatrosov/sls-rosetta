#!/bin/bash
# Build serverless functions as a Go plugin for the Yandex Cloud Functions runtime.
/golang/bin/go build -tags=ycf -buildmode=plugin -o /build/handler.so .

# if no /build/shared-libs exists, then create it
if [ ! -d "/build/shared-libs" ]; then
  mkdir -p /build/shared-libs
fi

# copy all shared libraries to /build/shared-libs
libs=$(ldd /build/handler.so | awk -F " => " '{split($2, a, " "); print a[1]}')
for l in $libs; do
  f="${l##*/}"
  cp "$l" "/build/shared-libs/$f"
done
