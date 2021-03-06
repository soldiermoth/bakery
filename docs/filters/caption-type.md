---
title: Caption Type
parent: Filters
nav_order: 4
---

# Caption Type
Values in this filter define a whitelist of the caption types you want to **EXCLUDE** in the modifed manifest. Passing an empty value for this filter will return all captions available in the manifest.

## Protocol Support

HLS | DASH |
:--:|:----:|
yes | yes  |

## Supported Values

| codec      | values | example  |
|:----------:|:------:|:--------:|
| Subtitles  | stpp   | ct(stpp) |
| WebVTT     | wvtt   | ct(wvtt) |


## Usage Example 
### Single value filter:

    $ http http://bakery.dev.cbsivideo.com/ct(stpp)/star_trek_discovery/S01/E01.m3u8

    $ http http://bakery.dev.cbsivideo.com/ct(wvtt)/star_trek_discovery/S01/E01.m3u8


### Multi value filter:
Mutli value filters are `,` with no space in between

    $ http http://bakery.dev.cbsivideo.com/ct(stpp,wvtt)/star_trek_discovery/S01/E01.m3u8

