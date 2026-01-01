FROM dadebia/codebox-base-docker:latest

COPY ./bin/ /codebox/bin
COPY ./docker/production.env /codebox/bin/codebox.env
COPY ./app/build /codebox/frontend
COPY ./docker/nginx.conf /etc/nginx/sites-enabled/codebox.conf
COPY ./docker/runserver.sh /codebox/
COPY ./atlas.hcl /codebox/bin
COPY ./migrations/ /codebox/bin/migrations
COPY ./templates/ /codebox/bin/templates

RUN chmod +x /codebox/runserver.sh && \
    chmod +x /codebox/bin/codebox

EXPOSE 8000

WORKDIR /codebox

VOLUME [ "/codebox/data" ]

CMD ["/codebox/runserver.sh"]
