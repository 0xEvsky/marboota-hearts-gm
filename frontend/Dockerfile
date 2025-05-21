# ~~~~~ Base ~~~~~
FROM node:23-alpine AS base

COPY . /app/frontend

WORKDIR /app/frontend

RUN npm install

EXPOSE 3000

CMD [ "npm", "run", "devnotunnel" ]
