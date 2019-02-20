docker-compose down -v
git pull
# docker rmi -f $(docker images -aq)
docker-compose build
docker-compose up ssp