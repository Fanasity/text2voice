FROM ubuntu
WORKDIR /service
COPY text2voice .
COPY conf ./conf
# VOLUME ["/service/conf"]
# 设定时区
ENV TZ=Asia/Shanghai
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
#RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
#RUN mkdir /lib64 && ln -s /lib/ld-musl-aarch64.so.1 /lib64/ld-musl-aarch64.so.2
RUN apt-get -qq update && apt-get -qq install -y --no-install-recommends ca-certificates curl
# ENTRYPOINT /service/text2voice
CMD /service/text2voice
