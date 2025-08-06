# ha2trmnl

This tool is the backend of a Private TRMNL plugin. It is very simple:

- Call the HomeAssistant entity API to get all entities
- Call the HomeAssistant websocket API to get all labels assigned in the UI
- Compute some metrics/aggregates based on entity states

Depending on the run mode, this aggregate is handled differently:

- `push` will run on a configurable interval, and send the data to a TRMNL webhook
- `serve` will expose a `/plugin_data` HTTP endpoint that exposes the data
- `fetch` will just log the data to STDOUT

You can run the tool via:

```bash
go run ./cmd/ha2trmnl [serve|push|fetch] (./path/to/config.yaml)
```

I added a `Dockerfile` that wraps this script to run it in a container as a CronJob.

Using this [markup.html](examples/markup.html), we get this nice plugin:

![screenshot](examples/plugin.png)

***This script is very customized to my needs and my HA setup. I expect the best use case for this repo will be an example and inspiration for others to build their own private plugins.***

This is an example of the data the tool will generate and give to TRMNL:

```json
{
  "merge_variables": {
    "generated": "17h13m",
    "lights": {
      "off": 31,
      "on": 32,
      "percent_on": 50.79365079365079
    },
    "metrics": {
      "humidity": {
        "garage": 50,
        "inside": 49.279999999999994,
        "outside": 59
      },
      "temperature": {
        "garage": 81,
        "inside": 68.86999999999999,
        "outside": 75.7
      }
    },
    "open_sensors": [
      "12d+ Mudroom Door",
      "5h+ Master Bathroom Window T"
    ],
    "speakers_playing":  [
      "Office"
    ]
  }
}
```

## Setup

First, make a new [Private Plugin](https://docs.usetrmnl.com/go/private-plugins/templates):

- Use either the Webhook or Polling strategy
  - For the Webhook option, copy the webhook URL (you need it for the config later)
  - For the Polling option, configure the following:
    - URL: `https://yoursite.com/plugin_data`
    - Headers: `authorization=bearer {{ api_key }}&content-type=application/json`
    - Form fields (to be able to set the `api_key`):
      ```yaml
      - keyname: api_key
        field_type: password
        name: API Key
        description: Bearer token passed to poller
        placeholder: s3cret!
      ```
- Edit the markup for the plugin. Here's an example: [markup.html](examples/markup.html)

Next, create config file, like what's in [example_config.yaml](examples/example_config.yaml)

Then, run this tool with docker compose:

```yaml
version: "3.8"

services:
  ha2trmnl:
    image: fisherevans/ha2trmnl:latest
    restart: unless-stopped
    command: ["serve", "/config/config.yaml"]
    ports:
      - "9123:9123" # this needs to match what's in your config
    volumes:
      - /path/to/your/ha2trmnl_config.yaml:/config/config.yaml:ro
```

## Developing

- `go run ./cmd/ha2trmnl serve` to run the tool locally based on the source code
- `./build.sh` will build a new docker image locally
- `./run_docker.sh push /config/config.yaml` will run the local docker image (useful for testing)
- `./publish.sh` will publish the docker image (only I can, haha)