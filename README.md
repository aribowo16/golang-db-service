# Golang DB Service

## Pre-requisite
1. Install golang v1.21 or above.
2. Install Postgres v12 or above.
  
## Postgres Table

```sql
CREATE TABLE users (
  userid SERIAL PRIMARY KEY,
  name TEXT,
  age INT,
  location TEXT
);
```

## How to Start
1. Create new database in Postgres
2. Run the SQL script above to create table users 
3. Open .env file inside the project directory
4. Change POSTGRES_URL value to your Postgres configuration<br>
```postgres://{username}:{password@{hostname}:{port}/{dbname}?sslmode=disable```
5. Open terminal 
6. Go to the project directory
7. Run command ```go run .```
8. Service will running on the port 8080

## API List
1. <b>Get User By ID</b><br>
Method: GET<br>
Path: ```/api/user/{id}```
2. <b>Get All Users</b><br>
Method: GET<br>
Path: ```/api/user```
3. <b>Insert User Data</b><br>
Method: POST<br>
Path: ```/api/user```<br>
Request Data: 
```json
{
    "id": 1, 
    "name": "Budi",
    "location": "Jakarta",
    "age": 20
}
```

## Author
Enrico Siswanto@NTT Data - 2023
