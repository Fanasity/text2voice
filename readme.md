# 将文本转成语音

通过kafka接受转语音的任务，将语音文件存储到minio中， 并通过kafka返回处理结果


### 构建镜像
docker build -t t2v .

### 运行docker容器
 docker run -d --net=host --name t2v -v /etc/localtime:/etc/localtime t2v