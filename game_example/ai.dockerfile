# 使用官方 Python 镜像作为基础镜像
FROM python:3.8-slim-buster

# 设置工作目录
WORKDIR /app

# 将 main.py 复制到工作目录中
COPY AI_PATH /app/ai.py

# 安装依赖
RUN pip install --no-cache-dir -r requirements.txt

# 暴露端口，使得应用可以被访问
EXPOSE 12345

# 定义环境变量
ENV NAME World

# 当容器启动时运行 Python 应用
CMD ["python", "ai.py"]