FROM python:3.11-alpine

WORKDIR /app

RUN pip install flask==2.2.3

# mount transcription folder as external volume
VOLUME /app/transcriptions

COPY ./api/static /app/static
COPY ./api/server.py .
COPY ./api/database.py .

CMD ["flask",  "--app", "server", "run", "--host=0.0.0.0"]