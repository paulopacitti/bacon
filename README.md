# bacon 🥓
A cli to update the public IP of your domain registered in [Porkbun](https://porkbun.com/). 

![a pig delivering mail](docs/pig.jpg)

This CLI is a wrapper around the [Porkbun API](https://porkbun.com/api/json/v3/documentation) to make it easier to update DNS with the machine's current public IP. Updating your domain registered to point to your Raspberry Pi would be a good example on why would you use this tool. Later, you could use a cronjob to call `bacon` periodically to update the your domain to point to your machine.

## Installing
Make sure you have `go` installed and that `$(go env GOPATH)/bin` is in your path, then:

```sh
go install github.com/paulopacitti/bacon@latest
```

## Usage
1. Add the keys and domain you want to update: `bacon config` 
2. Get your public IP and update the domain that you configured previously: `bacon update`

Better explanation of each command can be found on `bacon -h`. The configuration will be saved in `$HOME/.config/bacon/config.json`.

## Contributing
Feel free to open PRs to fix bugs and add new features 🚀  