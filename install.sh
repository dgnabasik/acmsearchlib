echo "Still need timeevent."
cd ~/go/src/github.com/dgnabasik/acmsearchlib
echo dgnabasik
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.0"
git push -u origin main

cd ./nulltime
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./headers
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./filesystem
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./database
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./timestampinterval
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./vocabulary
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./article
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./conditional
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./wordscore
 go get -u ./...
 go build
 go install
 echo -n "?"
 read 
 cd ..

go get -u github.com/dgnabasik/acmsearchlib
echo "done!"
