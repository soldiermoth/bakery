<p align="center">
  <img width="250" src="http://flv.io/bakery-logo.svg">
</p>

[![codecov](https://codecov.io/gh/cbsinteractive/bakery/branch/master/graph/badge.svg)](https://codecov.io/gh/cbsinteractive/bakery)


Bakery is a proxy and filter for HLS and DASH manifests.

## Setting up environment for development

#### Clone this project:

    $ git clone https://github.com/cbsinteractive/bakery.git

#### Export the environment variables:

    $ export BAKERY_CLIENT_TIMEOUT=5s 
    $ export BAKERY_HTTP_PORT=:8082
    $ export BAKERY_ORIGIN_HOST="https://streaming.cbs.com" 

Note that `BAKERY_ORIGIN_HOST` will be the base URL of your manifest files.

#### Run the API:

    $ make run

The API will be available on http://localhost[:BAKERY_HTTP_PORT]

## Run Tests

    $ make  test

## Help

You can find the source code for Bakery at GitHub:
[bakery][bakery]

[bakery]: https://github.com/cbsinteractive/bakery

If you have any questions regarding Bakery, please reach out in the [#i-vidtech-mediahub](slack://channel?team={cbs}&id={i-vidtech-mediahub}) channel.
