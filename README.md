# Go X509

Client server programs connecting with criptography using three-pass x.509 authentication protocol

## Getting Started

These instructions will get you a copy of the project up and running on your local machine for development, testing or using purposes.

### Prerequisites

* Go version 1.5 at least
* Some linux distro with make

### Installing

A step by step series of examples that tell you have to get a development env running

Clone the project

```bash
git clone git@github.com:rafaelescrich/go-x509.git
cd $GOPATH/src/github.com/rafaelescrich/go-x509
```

Build binary with make tool

```bash
make
```
Then if everything runned smoothly you should have a binary
To run it, just type

```bash
chmod +x makecert.sh
./makecert.sh joe@random.com
```

```bash
./client
```

```bash
./server
```

## Running the tests

make test

## Built With

* [Argon2](https://github.com/golang/crypto/tree/master/argon2) - Go supplementary cryptography libraries

## TODO

* Testing:

## Author

* **Rafael Escrich** - [github.com/rafaelescrich](https://github.com/rafaelescrich)

## License

This project is licensed under the GPL v2 License - see the [LICENSE.md](LICENSE.md) file for details