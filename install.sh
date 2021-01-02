cd ~/go/src/github.com/dgnabasik/acmsearchlib
echo dgnabasik
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.0"
git push -u origin main

cd ./nulltime
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./headers
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./filesystem
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./database
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./timestampinterval
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./vocabulary
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./conditional
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

cd ./wordscore
 go mod tidy
 go build
 go install
 echo -n "?"
 read 
 cd ..

go get -u github.com/dgnabasik/acmsearchlib
echo "done!"
