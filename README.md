# bacon 🥓
A cli to update the public IP of your domain registered in [Porkbun](https://porkbun.com/). 

![](docs/pig.jpg)

This CLI is a wrapper around the [Porkbun API](https://porkbun.com/api/json/v3/documentation) to make it easier to update DNS with the machine's current public IP. Updating your domain registered to point to your Raspberry Pi would be a good example on why would you use this tool. Later, you could use a cronjob to call `bacon` periodically to update the your domain to point to your machine.
## Installing

## Usage
1. Add the keys and domain you want to update. This configuration will be saved in `$HOME/.config/bacon/config.json`:
 
   `bacon config` 
2. Get your public IP and update the domain that you configured previously:
   
   `bacon update`

## Contributing
Feel free to open PRs to fix bugs and add new features 🚀  