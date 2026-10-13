# Turbo Jwt

## Description

TurboJwt is a simple librarie for jwt authentication for go app

## Installation

To add turbojwt in your projet use the following command :

```bash
go get github.com/MaminirinaEdwino/turbojwt
```

## Usage

1. To generate a token

```go
// import the turbojwt package
import "github.com/MaminirinaEdwino/turbojwt"

// In your main function or your logic 

payload := map[string]any{
    "user": User{
        Name: "username",
        Age: 24,
    },
}

secret := "my_secret_key"
token, _ := turbojwt.Encode(secret, payload)

//show the result
fmt.Println(token)
```

1. Verify a token

```go
// import the turbojwt package
import "github.com/MaminirinaEdwino/turbojwt"

// In your main function or your logic 
// Here token is string
fmt.Println(turbojwt.Verify(secret, token))
```
