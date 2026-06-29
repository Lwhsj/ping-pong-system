# 构建阶段：编译前端
FROM node:22 AS build
WORKDIR /app
# 先复制依赖清单，利用分层缓存
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

# 运行阶段：仅保留静态文件，由 nginx 提供服务并反向代理后端 API
FROM nginx:alpine
COPY nginx.conf /etc/nginx/conf.d/default.conf
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
