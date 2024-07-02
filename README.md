# DAR - Data At Rest

This is a small program that allows you to store data at rest in an encrypted format.

## Install

```sh
go install github.com/cthain/dar
```

## Usage

```
dar <mode> [options]
where <mode> is one of:
    -d | -decrypt: Decryption mode. The input data will be decrypted.
    -e | -encrypt: Encryption mode. The input data will be encrypted.

Options:
    -k <string> | -key <string>: The pass key to use to decrypt/encrypt the data.
                                 If the -key option is not provided you will be prompted to enter the key. It is not recommended to pass this option since it could expose
                                 your pass key in plaintext. It is meant to support tooling where the
                                 key is provided through an environment variable.
    -i <string> | -in <string>:  The name of the file to read input data from. By default the input is
                                 read from stdin.
    -o <string> | -out <string>: The name of the file to write output data to. By default the output is
                                 written to stdout.
```

## Examples

### Encrypt and decrypt some data

```sh
$ echo "this is plaintext" | ./dar -e -k "5up3r s3cReT" | ./dar -d -k "5up3r s3cReT"
this is plaintext
```

### Encrypt a plaintext file to an output file

```sh
# encrypt the input file data
$ dar -e -in plain.text -out encypted.dar
Enter pass key to encrypt: 
wrote encrypted data to encrypted.dar

# decrypt the encrypted file data
$ dar -d -in encrypted.dar
Enter pass key to encrypt: 
this is plaintext
```
