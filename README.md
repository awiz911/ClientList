# ClientList

Description 

This project is a Go learning project with simple RestFul API. 
It uses postgres db inside with docker-compose. You can compose with dockerfile or create your own postgres database without it.

To run docker-compose(in project folder): 

     cd docker
     docker-compose up
      
PostgreSQL works on 32300 Port (32300 -> 5432). You can access with database IDE (DataGrip etc.) with configure port 32300.

If you want to connect from your host system type the following command to terminal.

psql -h localhost -p 32300 -d docker -U docker --password

Database table configuration
      
      CREATE TABLE USERS (
        ID INT PRIMARY KEY,
        NAME TEXT NOT NULL,
        LASTNAME TEXT NOT NULL,
        AGE INT NOT NULL
        CELL TEXT NOT NULL,
        EMAIL TEXT NOT NULL
        );
      
      CREATE SEQUENCE public.users_id_seq NO MINVALUE NO MAXVALUE NO CYCLE;
      ALTER TABLE public.users ALTER COLUMN id SET DEFAULT nextval('public.users_id_seq');
      ALTER SEQUENCE public.users_id_seq OWNED BY public.users.id;

Database access configuration inside code
Under config/config.go directory in the project, you will find database access configuration. 
You can change it with your custom configuration(need to change it in docker-compose.yml as well)

      DB_USER     = "docker"
      DB_PASSWORD = "docker"
      DB_NAME     = "docker"
      PORT = "32770"

To start service(after starting PSQL in docker) type commands in your working terminal:

1. Get package neccessary for postgreesql in go.

     go get github.com/lib/pq 
     
2. Start service

      go run main.go

If there is no error then you should be able to work with service on http://127.0.0.1:3000

http://127.0.0.1:3000/clients/

1. List all clients (GET)

2. Add new client with JSON(POST)

      {
      	"name": "mockName",
      	"lastname": "mockSurname",
      	"age": 30
        "cell": 348094834
        "email": mockemail@yahoo.com
      	}

http://127.0.0.1:3000/clients/<ID>

1. List one client with given ID (GET)

2. Update one client with given ID (PUT)

3. Delete one client with given ID (DELETE)
