FROM node:15-alpine

WORKDIR /app

COPY . .

RUN npm install

RUN npm run build

CMD ["node", "server.js"]
