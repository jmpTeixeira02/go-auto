# GO-Auto

Go-Auto is a tool built in [Go](https://go.dev/) that scrapes [StandVirtual](https://www.standvirtual.com/) and notifies on new automotive deals.

### Database Options
Go-Auto is capable of supporting multiple storage options, currently supporting:
- [SQLite](https://www.sqlite.org/)

### Notifiers
Go-Auto is capable of supporting different types of notifiers, currently supporting:
- [Discord](https://discord.com/)
- Terminal

## How to Use
1. Go to [StandVirtual](https://www.standvirtual.com/), search with the desired filters and copy the URL
2. Open config/config.yml and edit the fields
```
    url: <URL from step 1>
    refresh_min: <how often to check in minutes>
    notifier:
        service: [terminal, discord]
        config:
            token: <notifier token, only needed with discord>
            receiver: <notifier receiver, only needed with discord>
    data:
        service: sqlite
        config:
            address: <address of the db> 
```

### Example Config
```
url: https://www.standvirtual.com/carros/desde-2014?search%5Bfilter_float_first_registration_year%3Ato%5D=2022&search%5Bfilter_float_mileage%3Ato%5D=10000&search%5Bfilter_float_price%3Ato%5D=20000&search%5Badvanced_search_expanded%5D=true
refresh_min: 30
notifier:
  service: terminal
  config: 
    token: "token"
    receiver: "reciver"
data:
  service: sqlite
  config:
   address: "tmp/sqlite.db"
```

## How to Run

### Binary
To build and run locally simply execute the following commands
```
make deps
make
./tmp/auto
```

### Docker
You can run the tool in [Docker](https://www.docker.com/) with [Compose](https://docs.docker.com/compose/) by running

        make docker
