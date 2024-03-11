# Pocket Counter

A web application to visualize how the number of unread items in Pocket changes in time.


## Setup

On the VPS, add this to `docker-compose.yml`:
```yml
  pocketcounter:
     image: pviotti/pocket-counter:latest
     container_name: pocket-counter
     restart: unless-stopped
     env_file:
      - .env
     volumes:
      - ./data/pocket-counter:/app/data
```

Then add appropriate environmens variables in `.env`.