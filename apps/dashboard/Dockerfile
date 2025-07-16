FROM oven/bun:1

WORKDIR /app

RUN bun install --global nuekit

COPY . .

EXPOSE 8081

CMD ["sh", "-c", "nue build -p && nue serve -p"]
