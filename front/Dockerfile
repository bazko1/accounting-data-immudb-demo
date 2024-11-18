FROM node:23-alpine3.19
WORKDIR /app

COPY . .
RUN npm install
RUN npm run build --production
RUN npm install -g serve
CMD serve -s build
