# locate work space
work_dir=$(pwd)

# change directory to work space
cd ${work_dir}

# clean old output by force
rm -rf ./output

# create new output
mkdir output
cd output
mkdir bin
mkdir data
mkdir conf

# copy data and conf
cd data
cp ../../data/url.data ./
cd ../conf
cp ../../config/spider.conf ./

cd ${work_dir}

# download dependency
go mod download

# go test
go test -v -gcflags=-l ./...

if [ $? -ne 0 ]
    then
        echo "go test failed"
        exit -1
fi

# build
go build -o ./output/bin/mini_spider

cd output/bin
chmod 777 mini_spider

echo "crawling ..."
# crawl
./mini_spider -c ../conf -l ../log