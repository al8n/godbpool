set -e
echo "" > coverage.txt

cd ./gormpool
./test.sh

cd ../
