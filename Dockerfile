FROM registry.cn-shenzhen.aliyuncs.com/xxxjz/alpine

# 使用ARG定义参数
ARG APP_NAME
ARG PROFILES

ENV TZ Asia/Shanghai
ENV CONFIG_FILE=./config/config.yaml
# 使用参数设置环境变量
ENV APP_NAME ${APP_NAME}
ENV PROFILES ${PROFILES}

WORKDIR /app

COPY ./bin/xxxjz .
COPY ./config/config_${PROFILES}.yaml $CONFIG_FILE
COPY ./docs ./docs

EXPOSE 8080

RUN chmod +x ./xxxjz
#CMD sh -c 'echo "构建参数值为: $PROFILES $APP_NAME"'
CMD sh -c './xxxjz $APP_NAME -c ./config'
