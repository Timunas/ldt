FROM bitnami/minideb:jessie

WORKDIR /app
COPY ./ldt-server /app/ldt-server
COPY ./.env .env

CMD ["/app/ldt-server"]
