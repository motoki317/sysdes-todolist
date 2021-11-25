FROM node:16-alpine AS build

WORKDIR /app

COPY ./package*.json .
RUN yarn install --frozen-lockfile && yarn cache clean

COPY . .
RUN yarn build

FROM caddy:latest

COPY --from=build /app/dist/spa /usr/share/caddy
