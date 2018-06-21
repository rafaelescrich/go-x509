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
git clone git@github.com:rafaelescrich/go-x509.git $GOPATH/src/github.com/rafaelescrich/go-x509
cd $GOPATH/src/github.com/rafaelescrich/go-x509
```

Build binary with make tool

```bash
make all
```

Create the directory of the certificates

```bash
mkdir certs
```

After we have a binary, we must first generate our keys

```bash
./go-x509 -g server
./go-x509 -g client
```

Then if everything runned smoothly you should have a binary
To run it in server mode, just type

```bash
./go-x509 -s -p 8000
```

Now we need to run the client

```bash
./go-x509 -c -p 8000
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