#!/bin/bash
# Example use:
# TAG=v1.2 TOKEN=447879c5af5887eab22725605783e86d3304bc99 ./.github_release.sh

if [ -z "$TAG" ]; then
  echo "TAG is empty please do 'export TAG=v0.1' for example"
  exit 1
fi

if [ -z "$TOKEN" ]; then
  echo "You forgot write yout TOKEN"
  exit 1
fi

make || exit 1
make release

LINUX_BIN_PATH=builds/logsend
DARWIN_BIN_PATH=builds/logsend_darwin

gzip -9 -f $LINUX_BIN_PATH || exit 1
gzip -9 -f $DARWIN_BIN_PATH || exit 1


response=`curl --data "{\\"tag_name\\": \\"$TAG\\",\\"target_commitish\\": \\"master\\",\\"name\\": \\"$TAG\\",\\"body\\": \\"Release of version $TAG\\", \\"draft\\": false,\\"prerelease\\": false}" \
  -H 'Accept-Encoding: gzip,deflate' --compressed "https://api.github.com/repos/ezotrank/logsend/releases?access_token=$TOKEN" > response`

release_id=`cat response|head -n 10|grep '"id"'|head -n 1|awk '{print $2}'|sed -e 's/,//'`
rm response

if [ -z "$release_id" ]; then
  echo "something wrong"
  echo $response
  exit 1
fi

curl -X POST -H 'Content-Type: application/x-gzip' --data-binary @$LINUX_BIN_PATH.gz "https://uploads.github.com/repos/ezotrank/logsend/releases/$release_id/assets?name=logsend.gz&access_token=$TOKEN"
curl -X POST -H 'Content-Type: application/x-gzip' --data-binary @$DARWIN_BIN_PATH.gz "https://uploads.github.com/repos/ezotrank/logsend/releases/$release_id/assets?name=logsend_darwin.gz&access_token=$TOKEN"
