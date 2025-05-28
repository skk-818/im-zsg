edge:
    1、维护 ws 的长连接
    2、数据协议解析
    3、维护会话
    4、推送消息到 客户端

# 查看 topic
docker exec -it kafka kafka-topics.sh --bootstrap-server localhost:9092 --list
# 创建 topic
docker exec -it kafka kafka-topics.sh --create --topic topic-edge-01 --bootstrap-server localhost:9092 --partitions 3 --replication-factor 1