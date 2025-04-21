# FLKCLI

FLKCLI is a command-line interface (CLI) tool that allows you to interact with the Flickr API. With FLKCLI, you can perform various operations such as managing photo sets.

## Features

- Login to flickr
- List photo sets
- Upload photos

## Installation

To install FLKCLI, follow these steps:

1. Clone the repository: `git clone https://github.com/Sylbat92/flkcli.git`
2. Navigate to the project directory: `cd flkcli`
3. Build: `make`
4. Make executable and add to PATH `chmod +X flkcli && cp flkcli /usr/local/bin/`

## Usage

To use FLKCLI, follow these steps:

1. Obtain a Flickr API key by creating a new app on the [Flickr Developer Portal](https://www.flickr.com/services/apps/create/).
2. Configure FLKCLI with your API key: `flkcli setup --apikey YOUR_API_KEY --apisecret YOUR_API_SECRET`
3. Login to your Flickr account : `flkcli login`