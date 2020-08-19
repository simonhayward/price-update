# price-update

For a set of securities identitifed by their ISIN update their current prices.


- From external files populate securites and store objects.
- For each ISIN create a task for each its sources to retrieve the price. When a task has succesfully found a price cancel its other tasks.
- Add/Update the store ISINs with the new prices.
- Save the store back to its external location.


## Required

#### Environment variables

- `SOURCES` - URL to populate `Securities`
- `STORE` - URL to populate `Store`
- `STORE_UPDATE` - URL to patch changes
- `STORE_TOKEN` - authorization token for patch request

## External files

#### Sources format (JSON)


```
{
    "isin": <International Securities Identification Number>,
    "src": [
        {
            "url": <Web resource>,
            "p": <Regex pattern>
        }
    ]
}
```

```
[
    {
        "isin": "GB000434554312",
        "src": [
            {
                "url": "http://localhost:8000/1.html",
                "p": "Price <\/span><span class=\"price\">([.0-9]+)<\/span>"
            },
            {
                "url": "http://localhost:8000/2.html",
                "p": "sell price<\/div><h3>([.0-9]+)p<\/h3>"
            }
        ]
    }
]

```

#### Store format (CSV)

```
price,isin,updated
```


```
320.21,GB000434554312,2006-01-02T15:04:05-0700
578.70,GB0005454854356,2006-01-02T15:04:05-0700
```

## Running

#### Binary

Executing the binary having exporting the envs.

```
go build && ./price-update
```

#### Docker

```
docker build -t price-update .
docker run --env SOURCES=$SOURCES --env STORE=$STORE \
        --env STORE_UPDATE=$STORE_UPDATE --env STORE_TOKEN=$STORE_TOKEN \
        --rm price-update
```

#### Serverless

This currently deploys out as [GCP function](https://cloud.google.com/functions/) which is then run by a [scheduler](https://cloud.google.com/scheduler/).

```
make deploy
```
