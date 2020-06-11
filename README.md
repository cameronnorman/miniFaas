# MiniFaas
MiniFaas is a small golang application that allows you to create and run binary files through an API. You can also run system commands such as `ls` or `cat` if you want to interact with your server

## Getting Started
To run this application you can use docker. Start the container by running the following command:

### Build Docker image:
```
docker build . -t faas:1.0.0
```

### Run with Docker:
```
docker run --rm --env TOKEN=<your-secure-token> -p <your-desired-port>:1312 -d faas:1.0.0
```

## How to use

### Creating functions
To create a function that can be used through the API you can make the following HTTP request

```
POST 'http://localhost:1312/create?TOKEN=<your-secure-token>'

Headers:
  Content-Type: "multipart/form-data"
Body:
  file: (uploaded-binary-file) (file: say_hello)
  name: the-name-of-your-function (name: say_hello)
```

```
output:

{
    "message": "Upload successful"
}
```

### Listing functions
To list the functions that have been created you can make the following request

```
GET 'http://localhost:1312/all?TOKEN=<your-secure-token>'
```

### Running functions
To run a created function you can use the following API request

```
POST 'http://localhost:1312/say_hello/run?TOKEN=<your-secure-token>'

Headers:
  Content-Type: "application/json"
Body:
  name: "John"
```

```
output:

{"response": "Hello John"}
```

If there is an error while interacting with your binary you are likely to get the following response:
```
output:

{"error": "There was an error trying to say hello"}
```

## Todos:
- [ ] Write tests
