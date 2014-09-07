#!/bin/bash

function release(){
  if [ -z "$TAG" ]; then
    echo "TAG is empty please do 'export TAG=v0.1' for example"
    exit 1
  fi

  if [ -z "$TOKEN" ]; then
    echo "You forgot write yout TOKEN"
    exit 1
  fi

  make || exit 1

  gzip -9 -f -k $GOPATH/bin/logsend_linux || exit 1
  gzip -9 -f -k $GOPATH/bin/logsend || exit 1


  response=`curl --data "{\\"tag_name\\": \\"$TAG\\",\\"target_commitish\\": \\"master\\",\\"name\\": \\"$TAG\\",\\"body\\": \\"Release of version $TAG\\", \\"draft\\": false,\\"prerelease\\": false}" \
    -H 'Accept-Encoding: gzip,deflate' --compressed "https://api.github.com/repos/ezotrank/logsend/releases?access_token=$TOKEN" > response`

  release_id=`cat response|head -n 10|grep '"id"'|head -n 1|awk '{print $2}'|sed -e 's/,//'`
  rm response

  if [ -z "$release_id" ]; then
    echo "something wrong"
    echo $response
    exit 1
  fi

  curl -X POST -H 'Content-Type: application/x-gzip' --data-binary @$GOPATH/bin/logsend_linux.gz "https://uploads.github.com/repos/ezotrank/logsend/releases/$release_id/assets?name=logsend_linux.gz&access_token=$TOKEN"
  curl -X POST -H 'Content-Type: application/x-gzip' --data-binary @$GOPATH/bin/logsend.gz "https://uploads.github.com/repos/ezotrank/logsend/releases/$release_id/assets?name=logsend_darwin.gz&access_token=$TOKEN"

  rm -f $GOPATH/bin/logsend.gz $GOPATH/bin/logsend_linux.gz
}

while getopts "h?r --long release" opt; do
    case "$opt" in
    h|\?)
        show_help
        exit 0
        ;;
    r|release)  release
        ;;
    esac
done

