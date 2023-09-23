# auth-server
GRPC Go Auth Server

### Setup
Add a .env with the following variables:
- MONGO_URI: Your MongoDB SRV Connection String, e.g. mongodb+srv://[username]:[password]@[hostname.domain.TLD]/
- EMAIL_ADDRESS: For sending emails
- EMAIL_PASSWORD: Set up 2FA for Gmail then generate and add the app password
- EMAIL_HOST: smtp.gmail.com if you're using Gmail
- EMAIL_HOST_PORT: 465 if you're using Gmail

(Optional) If you made changes to the proto files, regenerate the files with the following command
```
make generate_grpc_code
```

### Build
```
go build
```

### Run
```
go run main.go
```
