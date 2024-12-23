# Event Planner

## Run in Docker

Clone git repo
```bash
git clone https://github.com/madalinpopa/go-event-planner.git && cd go-event-planner
```
Build docker image
```bash
docker build -t go-event-planner . 
```
Run docker container
```bash
docker run -p 4000:4000 \
    -v $(pwd)/database:/app/database \
    go-event-planner
```

