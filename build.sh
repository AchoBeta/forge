# build命令
RUN_NAME="achobeta.server.forge"
mkdir -p output/bin couput/conf
sp script/* output/
cp -r conf/* output/conf/
chmod +x output/bootstrap.sh

go build -o output/bin/${RUN_NAME} -coverpck=./...