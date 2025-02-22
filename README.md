# DISCLAIMER

This is not mean't to be used in production. This was built purely for educational purposes only. Use in your own environment at your own risk.

Currently I don't plan on creating this a finished project, but was created to brush up on Go and networking skills.

If you want a local DNS resovlver so you can have custom IP addresses I suggest you look through them listed here:

- [Awesome Selfhosted - DNS](https://github.com/awesome-selfhosted/awesome-selfhosted?tab=readme-ov-file#dns)

# About

Welcome to Simple-DNS a simple lightweight DNS resolver that allows you to create simple zone files using a yaml syntax.

## How it Works

The DNS resolver listens to the standard port `53` and will listen to any incoming DNS requests and tries to resolve them. The general process will look like below:

1. A query will be recieved on port 53
2. The host name will be checked against the current in-memory cache
3. Return the cached value if it exists
4. Check if the domain name is in loaded configuration and search the records
5. Return if record is found
6. Send the query to the next DNS resolver as defined in `config.yml`

## Features

- Caching
- YAML based configuration
- `A`, `CNAME`, `TXT` record storage

## Work in Progress

- [ ] Logging (yes I know I've just been lazy)
- [ ] Configuration file reloading without restarting
- [ ] Docker image
