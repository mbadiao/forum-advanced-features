docker build -t forum .
echo "...........docker images.........................."
docker images
echo "....................................."
docker container run -p  8081:8080 --detach --name containersforums forum 
echo "...........docker ps -a.........................."
docker ps -a
echo "....................................."
docker exec -it containersforums /bin/bash
echo "....................................."
ls -l