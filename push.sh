echo dgnabasik
echo aaf61ad438e564e493589b929e74437baba79df5
git pull https://github.com/dgnabasik/acmsearchlib
echo -n "push?"
read
git add --all :/
git commit -am "Release 1.0.11"
git push -u origin main
