# 编译

拉取最新tag代码

make build


go 版本>= 1.12

# 安装
https://github.com/bitcv/bacchain_node_linux



#DOCKER 运行方式
#创建bacd镜像

docker build -f Dockerfile_bacd   ./ -t bacd:1.0

#创建baccli镜像

docker build -f Dockerfile_baccli  .  -t baccli:1.0

#初始化运行环境

./run_init.sh

#运行bacd节点
docker run -it   -p26657:26657   -v  ~/.bacd:/root/.bacd     bacd:1.0  bacd start

#运行baccli
docker run -it   -p1317:1317   -v   ~/.baccli:/root/.baccli  baccli:1.0   baccli    rest-server --node={$local_ip}:26657 --chain-id=bacchain-mainnet-1.0 --laddr=tcp://0.0.0.0:1317