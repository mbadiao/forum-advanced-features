 # Stage 1: Build
FROM golang:alpine

# Metadata and labels
LABEL key="Forum"
LABEL maintainer="mbadiao, emalo, ousmasene, babacandiaye"
LABEL version="1.0"
LABEL description="A web forum for communication with posts and comments"

# Set the working directory to /app in the container
WORKDIR /app

# Copy all contents from the local build context into /app in the container
COPY . /app

# Utilisez 'apk' pour installer GCC, les bibliothèques de développement standard, Bash et SQLite3
RUN apk update && \
    apk add gcc libc-dev bash sqlite-dev sqlite

# Exécuter la commande 'go mod download' pour télécharger les dépendances du module Go
RUN go mod download

RUN CGO_ENABLED=1 GOOS=linux go build -o forum

EXPOSE 8080
# Exécuter la commande 'go build' pour compiler l'application, en créant un exécutable nommé "forum" dans le répertoire de travail

# Commande par défaut à exécuter lorsque le conteneur est démarré
CMD go run .
