# price-update

For a set of securities identitifed by their ISIN update their current prices.


- From external files populate securites (json) and output (csv).
- For each ISIN create a task for each its sources to retrieve the price. When a task has succesfully found a price cancel its other tasks.
- Add/Update the store ISINs with the new prices.
- Save the output back to its external location.


## Required

#### Environment variables

- `INPUT` - URL to populate `Securities`
- `OUTPUT` - URL to populate `Store`
- `API` - URL to patch changes
- `TOKEN` - authorization token for patch request

## External files

#### input.json


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

#### output.csv

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
docker run --env INPUT=$INPUT --env OUTPUT=$OUTPUT \
        --env API=$API --env TOKEN=$TOKEN \
        --rm price-update
```

#### Serverless

This currently deploys out as [GCP function](https://cloud.google.com/functions/) which is then run by a [scheduler](https://cloud.google.com/scheduler/).

```
make deploy
```
