FROM python:3.11-slim

RUN pip install sqlite-web

WORKDIR /data

ENTRYPOINT ["sqlite_web"]
CMD ["--host", "0.0.0.0", "--port", "8080", "mine.db"]

EXPOSE 8080
