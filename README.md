# Realtime Random Chat

## Server Side
Server side uses Go and Gorrilla

## Client Side
There are two client side applications you can choose, either Reactjs or Nodejs with express

# Running Server

```shell
$ make build

$ ./rumpi
```

# Running Client

## Node Client
https://github.com/musobarlab/rumpi/tree/master/client/node-client
```shell
$ cd client/node-client
$ npm install
$ npm start
```

Or with specific `PORT`

```shell
$ PORT=3000 npm start
```

Visit http://localhost:3000

## Reactjs Client
https://github.com/musobarlab/rumpi/tree/master/client/react-client
```shell
$ cd client/react-client
$ npm install
$ npm start
```