echo dgnabasik
git pull https://github.com/dgnabasik/acmsearchlib
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.12"
git push -u origin main
