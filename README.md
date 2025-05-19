# plsort

plsort is a tool that sorts YouTube playlist items in lexicographic order by title.

> [!note]
> Need to create a GCP project to use [YouTube Data API](https://developers.google.com/youtube/v3).
> 
> see https://developers.google.com/youtube/v3/getting-started

## Usage

```
Usage: plsort <playlist-id> [flags]

Arguments:
  <playlist-id>    YouTube playlist ID.

Flags:
  -h, --help                          Show help.
      --creds="client_secret.json"    OAuth client ID credentials JSON path.
  -r, --reverse                       Sort in reverse order.
      --version
```

## Getting Started

1. Create OAuth client ID, download the credentials and save them as `client_secret.json`.
    * see ｃhttps://developers.google.com/youtube/v3/quickstart/go
2. Get YoutTube Playlist ID.
    * ![](https://github.com/user-attachments/assets/54150b84-7d4a-4656-83ba-61a9314f1c2c)
3. Run plsort.
    ```sh
    $ plsort PLXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXXX
    ```
