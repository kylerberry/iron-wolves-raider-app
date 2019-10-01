**IN PROGRESS**

This is an app to aggregate Iron Wolves clan member stats by raid.

**How it works**

A service (cron) requests data from bungie every 15 minutes and then caches that response. Web requests will use the cache data. This is to avoid a very long wait time since it's roughly ~400 requests that need to be made to gather all the relevant data.

This was built using Golang as a fun introduction to the language.

**To test**

1) execute the server binary
`./server_exec`

2) get a raid leaderboard data dump:

- `http://localhost:8000/api/raid/lastwish/leaderboard`
- `http://localhost:8000/api/raid/crownofsorrows/leaderboard`
- `http://localhost:8000/api/raid/scourgeofthepast/leaderboard`

**TODO**
- cache bungie response data via cron (for quicker data fetching)
- active requests use cached response data
- UI for displaying the data
