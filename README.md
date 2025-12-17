A lightweight script to scrape video links from YouTube Shorts and TikTok based on a specific date range.
> Instagram support comming soon

<hr>

### Setup & Configuration

1. **Clone the repo**
2. **Configure credentials (refer to [config.example.json](https://github.com/llan0/link-fetcher/blob/main/config.example.json))**
3. **Fill in credentials** 
> YouTube: Requires Google Cloud API Key with "YouTube Data API v3" enabled <br>
> TikTok: Requires a RapidAPI key (specifically [tiktok-api23](https://rapidapi.com/Lundehund/api/tiktok-api23/))

<hr>

### Usage

Refer to the included [Makefile](https://github.com/llan0/link-fetcher/blob/main/Makefile). Results are saved to the `output/` directory 

```
make run-all
make tiktok 
make youtube
make instagram (WIP)
```

<hr>

### Project Structure

```
.
├── config.example.json   
├── config.json           # local config 
├── go.mod
├── go.sum
├── internal
│   ├── config
│   │   └── config.go    
│   ├── fetcher
│   │   ├── instagram.go  # (WIP)
│   │   ├── platform.go   
│   │   ├── tiktok.go     # TikTok implementation
│   │   └── youtube.go    # YouTube implementation
│   └── writer
│       └── writer.go     # unified ouptput writer 
├── main.go              
├── Makefile              
└── README.md
```

<hr>

### TODO
[ ] Add Instagram support <br>
[ ] Better error handling  <br>
[ ] Custom output format <br>
[ ] Unify platform logic <br>
[ ] Add tests <br>
