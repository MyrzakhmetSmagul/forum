build:
    docker build -t forum .

run:
    docker run -d --name forum-runer -p 8080:8080 forum