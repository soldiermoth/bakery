<p align="center">
  <img width="250" src="http://flv.io/bakery-logo.svg">
</p>


Bakery is a proxy and filter for HLS and DASH manifests.

## Setting up environment for development

#### Clone this project:

    $ git clone https://github.com/cbsinteractive/bakery.git

#### Export the environment variables:

    $ export BAKERY_CLIENT_TIMEOUT=5s 
    $ export BAKERY_HTTP_PORT=:8082
    $ export BAKERY_ORIGIN_HOST="https://streaming.cbs.com" 

Note that `BAKERY_ORIGIN_HOST` will be the base URL of your manifest files.

##### Usage Example:

If your playback URL is `http://streaming.cbsi.video/star_trek_discovery/S01/E01.m3u8`, your `BAKERY_ORIGIN_HOST` should be set to `http://streaming.cbsi.video`. The playback URL on the proxy will be `http://bakery.host.here/star_trek_discovery/S01/E01.m3u8`. 

If you want to apply filters, they should be placed right after the Bakery host. Following the example above, if you want to filter out all the levels that are outside of a given bitrate range, the playback URL should be: `http://bakery.host.here/b(1000,4000)/star_trek_discovery/S01/E01.m3u8`, where 1000Kbps and 4000Kbps are the lower and higher boundaries.

#### Supported Filters

- HLS

| name | example | description |
|------|---------|-------------|
| bandwidth | b(200,800) | An inclusive range of variant bitrates<br> to include in the modified manifest,<br> variants outside this range will be filtered<br> out. If a single value is provided, it will<br>define the minimum bitrate desired in the <br>modified manifest


- DASH

| name | example | description |
|------|---------|-------------|
| stream | fs(audio) | Values in this filter define stream types you wish<br> to remove from your manifest. The filter in this  <br>example will filter out all audio streams from the<br> modified manifest. |
| caption type | ct(wvtt) | Values in this filter define a whitelist of the caption<br> types you want included in the modifed manifest.<br> Passing an empty value fopr this filter will remove<br>all caption types from the manifest.

#### Run the API:

    $ make run

The API will be available on http://localhost[:BAKERY_HTTP_PORT]

## Run Tests

    $ make  test
