# tsignal
Analyzing social media sentiment and its impact on stock market

> Personal project, does not represent Google

![tsignal topology](/../master/images/tsignal-topology.png?raw=true "tsignal topology")

This Go app deploys into [GKE](https://cloud.google.com/container-engine/), subscribes to Twitter stream for all companies defined in the `Stocks` table in [Cloud Spanner](https://cloud.google.com/spanner/) and scores each event against the [Google NLP API](https://cloud.google.com/natural-language/) while comparing the user sentiment against the stock ask price against Yahoo API. 

Eventually there will be a UI, for now, there are sample SQL scripts you can use to execute against the Cloud Spanner DB to analyze the data. 

> All GCP services used in this example can be run under the GCP Free Tier plan. More more information see https://cloud.google.com/free/

## Configuration

Edit the `scripts/config.sh` file with your Twitter API info. Alternatively
define the following environment variables. Instructions how to configure your Twitter API access codes are found [here](http://docs.inboundnow.com/guide/create-twitter-application/):

```
# export T_CONSUMER_KEY=""
# export T_CONSUMER_SECRET=""
# export T_ACCESS_TOKEN=""
# export T_ACCESS_SECRET=""
```

## Run

Once all the necessary environment variables are defined you can execute the `tsignal` appl locally `go run *.go` or use the included `Dockerfile` to create an image which you can then publish to the GCP (if needed, there is a script `build-publish` that will do all that for you.

## Cleanup

The cleanup of all the resources created in this application can be accomplished by executing the `scripts/cleanup.sh` script.

### TODO

* Tests, yes please
* UI for reports and config
* Way to subscribe to non-NASDAQ stock prices 
