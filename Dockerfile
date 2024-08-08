# 使用官方的 Node.js 作为基础镜像
FROM node:18
# 设置工作目录
WORKDIR /usr/src/app

# 复制 package.json 和 package-lock.json 到工作目录
COPY package*.json ./

# 安装项目依赖
RUN npm install

# 复制项目的所有文件到工作目录
COPY . .

EXPOSE 3001

# 定义环境变量（根据需要进行修改）
ENV NODE_ENV=production

# 启动应用程序
CMD ["node", "server.js"]
