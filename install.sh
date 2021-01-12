echo "Still need timeevent."
cd ~/go/src/github.com/dgnabasik/acmsearchlib

cd ./nulltime
 go get -u ./...
 go build
 go install
 echo -n "nulltime?"
 read 
 cd ..

cd ./headers
 go get -u ./...
 go build
 go install
 echo -n "headers?"
 read 
 cd ..

cd ./filesystem
 go get -u ./...
 go build
 go install
 echo -n "filesystem?"
 read 
 cd ..

cd ./database
 go get -u ./...
 go build
 go install
 echo -n "database?"
 read 
 cd ..

cd ./timestampinterval
 go get -u ./...
 go build
 go install
 echo -n "timestampinterval?"
 read 
 cd ..

cd ./conditional
 go get -u ./...
 go build
 go install
 echo -n "conditional?"
 read 
 cd ..

cd ./vocabulary
 go get -u ./...
 go build
 go install
 echo -n "vocabulary?"
 read 
 cd ..

cd ./article
 go get -u ./...
 go build
 go install
 echo -n "article?"
 read 
 cd ..

cd ./wordscore
 go get -u ./...
 go build
 go install
 echo -n "wordscore?"
 read 
 cd ..

echo dgnabasik
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.0"
git push -u origin main

go get -u github.com/dgnabasik/acmsearchlib
echo "done!"
