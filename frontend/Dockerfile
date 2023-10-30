FROM node:20-bullseye as builder
WORKDIR /app
COPY package.json ./
RUN yarn install
RUN yarn add serve
COPY . .
RUN yarn run build
CMD yarn serve -s build
