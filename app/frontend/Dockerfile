FROM node:18-alpine

ENV NODE_ENV=production
WORKDIR /app

COPY ["package.json", "package-lock.json*", "./"]

RUN npm install --production

COPY ["index.js", "instrument.js", "run.sh", "./"]

EXPOSE 4000

CMD [ "./run.sh" ]
