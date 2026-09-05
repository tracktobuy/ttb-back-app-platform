# ttb-back-app-platform
Backend service for our app platform


## Environment Variables

| Variable | Description | Required |
| --- | --- | --- |
| `MONGO_HOST` | Hostname (and connection options) of the MongoDB cluster | Yes|
| `MONGO_DB_USER` | Username used to authenticate with the MongoDB database |Yes |
| `MONGO_DB_PASSWORD` | Password used to authenticate with the MongoDB database | Yes|
| `MONGO_DB_NAME` | Name of the MongoDB database to connect to | Yes |
| `API_SERVER_PORT` | Port on which the API server listens. Default: 8080 | No |


## Run Locally with Docker

### First build the image

```
docker build -t ttb-backend-api .
```

### Run the container
```
docker run --name ttb-back-service \
-e MONGO_HOST=mycluster.subdomain.mongodb.net/?appName=mycluster \
-e MONGO_DB_USER=my-mong-user \
-e MONGO_DB_PASSWORD=super-secret--pss \
-e MONGO_DB_NAME=mongo-db-name \
-e API_SERVER_PORT=9999 \
-p 8080:9999 \
ttb-backend-api
```

### Database
Edit the ```.env-sample``` file accordingly to your database credentials data and then name it for ```.env```.