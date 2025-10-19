FROM dadebia/codebox-base-docker:latest

COPY ./bin/ /codebox/bin
COPY ./docker/production.env /codebox/bin/codebox.env
COPY ./app/build /codebox/frontend
COPY ./docker/nginx.conf /etc/nginx/conf.d/default.conf
COPY ./docker/startgin.sh /docker-entrypoint.d/
COPY ./atlas.hcl /codebox/bin
COPY ./migrations/ /codebox/bin/migrations
COPY ./html/ /codebox/bin/html

RUN chmod +x /docker-entrypoint.d/startgin.sh && \
    chmod +x /codebox/bin/codebox

EXPOSE 8000

WORKDIR /codebox

VOLUME [ "/codebox/data" ]

# go install github.com/swaggo/swag/cmd/swag@latest
# go install github.com/swaggo/swag/cmd/swag@latest
