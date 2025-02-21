FROM node:20-alpine3.20 AS build

COPY . /frontend

WORKDIR /frontend

RUN npm ci

RUN npm run build


FROM node:20-alpine3.20

COPY --from=build /frontend/build /frontend/build
COPY --from=build /frontend/package.json /frontend/package.json
COPY --from=build /frontend/package-lock.json /frontend/package-lock.json

WORKDIR /frontend

RUN npm ci --omit dev

EXPOSE 3000

CMD [ "node", "build" ]