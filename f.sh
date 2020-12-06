echo "go.mod"
grep -i $1 ./go.mod
echo "./filesystem/go.mod"
grep -i $1 ./filesystem/go.mod
echo "./database/go.mod"
grep -i $1 ./database/go.mod
echo "./timestampinterval/go.mod"
grep -i $1 ./timestampinterval/go.mod
echo "./headers/go.mod"
grep -i $1 ./headers/go.mod
echo "./nulltime/go.mod"
grep -i $1 ./nulltime/go.mod

