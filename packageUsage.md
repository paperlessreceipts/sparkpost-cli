# Momentum/SparkPost CLI Tools

[Source Code](https://github.com/SparkPost/sparkpost-cli)

The following CLI command use the RESTful API to interact with Momentum and SparkPost.

## Environment

All the CLI commands will check `SPARKPOST_APIKEY` and `SPARKPOST_BASEURL`.

If you are using the CLI command against SparkPost.com you do not need to set `SPARKPOST_BASEURL `.

* export SPARKPOST_API_KEY="VALID API KEY"
	* or use command line argument `--apikey "VALID API KEY"`
* export SPARKPOST_API_KEY="http://YOURSERVER.com"
	* or use command line argument `--baseurl "http://YOURSERVER.com"`

## Webhook CLI

The web hook CLI is a wrapper around the [Webhooks API](https://developers.sparkpost.com/api/#/reference/webhooks) which The the means to list, review and query web hook.

### Webhook List

[see](https://developers.sparkpost.com/api/#/reference/webhooks/list/list-all-webhooks)

List currently extant webhooks.

| Compatibility  | Compatible?  |
|---|:-:|
|SparkPost| Yes  |
|SparkPost Elite| Yes  |
|Momentum| Yes  |


* `> ./sp-webhook-cli`
* `> ./sp-webhook-cli --command list`

**Sample Output**

```
Name: “Delivery WebHook"
	hook ID:   5f61f8a0-738c-11e5-9579-0b90e3e7e87c
	Target:    http://webhook.domain.com:8080/xyz123
	Success:   2016-02-24T22:23:00+00:00
	Fail:      2016-02-24T21:53:00+00:00
	AuthType:  basic
```

### Query Webhook

[see](https://developers.sparkpost.com/api/#/reference/webhooks/retrieve/retrieve-webhook-details)

Retrieve details about a webhook by specifying its id in the URI path.


* `> ./sp-webhook-cli --command query --id 5f61f8a0-738c-11e5-9579-0b90e3e7e87c`

**Sample Output**

```
Name: "Yepher WebHook"
	hook ID:   5f61f8a0-738c-11e5-9579-0b90e3e7e87c
	Target:    http://webhook.domain.com:8080/xyz123
	Success:   2016-02-24T22:23:00+00:00
	Fail:      2016-02-24T21:53:00+00:00
	AuthType:  basic
	Events:
		bounce
		delivery
		injection
		spam_complaint
		out_of_band
		policy_rejection
		delay
		click
		open
		generation_failure
		generation_rejection
		list_unsubscribe
		link_unsubscribe
		relay_injection
		relay_rejection
		relay_delivery
		relay_tempfail
		relay_permfail
```


### Webhook Status

Retrieve details about a webhook by specifying its id in the URI path.

[see](https://developers.sparkpost.com/api/#/reference/webhooks/validate/retrieve-status-information)

Retrieve status information regarding batches that have been generated for the given webhook by specifying its id in the URI path. Status information includes the successes of batches that previously failed to reach the webhook's target URL and batches that are currently in a failed state.


* `./sp-webhook-cli --command status --id 5f61f8a0-738c-11e5-9579-0b90e3e7e87c`

**Sample Output**

```
BatchId: "24d44870-db40-11e5-b1e3-63a3a57c2125"
	Time:       2016-02-24T22:23:05.000Z
	Attempts:   4
	RespCode:   200
```



### Webhook CLI Help


```
--baseurl, -u "https://api.sparkpost.com"	Optional baseUrl for SparkPost. [$SPARKPOST_BASEURL]
--apikey, -k 				Required SparkPost API key [$SPARKPOST_API_KEY]
--username 				Username this is a special case it is more common to use apices
--password, -p 			Username this is a special it is more common to use apices
--verbose "false"		Dumps additional information to console
--command, -c "list"	Optional one of list, query, status. Default is "list"
--timezone, --tz 		Optional Standard timezone identification string, defaults to UTC Example: America/New_York.
--id 					Optional UUID identifying a web hook Example: 12affc24-f183-11e3-9234-3c15c2c818c2.
--limit 				Optional Maximum number of results to return. Defaults to 1000. Example: 1000.
--help, -h				show help
--version, -v			print the version

```

## Deliverability Metrics

This CLI is a wrapper around [Deliverability Metrics](https://developers.sparkpost.com/api/#/reference/metrics)

SparkPost and SparkPost (Elite) log copious amounts of statistical, real-time data about message processing, message disposition, and campaign performance. This reporting data is available in the UI or through the Metrics API. The Metrics API provides a variety of endpoints enabling you to retrieve a summary of the data, data grouped by a specific qualifier, or data by event type. Within each endpoint, you can also apply various filters to drill down to the data for your specific reporting needs.

| Compatibility  | Compatible?  |
|---|:-:|
|SparkPost| Yes  |
|SparkPost Elite| Yes  |
|Momentum| Yes  |

### Deliverability by Domain

Provides aggregate metrics grouped by domain over the time window specified. Use `--metrics` and `--domains` to control what columns or domains are returned. It is sometimes useful to pipe the output to a cvs file and open with Excel.

* `./sp-deliverability-metrics-cli --from "2014-02-01T00:00"`
* `./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "domain"`

### Deliverability Metrics by Binding

[see](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-binding/deliverability-metrics-by-binding)

**Note:** This endpoint is available in SparkPost Elite or Momentum only.

Provides aggregate metrics grouped by binding over the time window specified.

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "binding"``


### Deliverability Metrics by Binding Group

[see](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-binding)

**Note:** This endpoint is available in SparkPost Elite or Momentum only.Provides aggregate metrics grouped by binding group over the time window specified.

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "binding-group"``


### Deliverability Metrics by Campaign

[see](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-binding/deliverability-metrics-by-campaign)

Provides aggregate metrics grouped by campaign over the time window specified.

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "campaign"``


### Deliverability Metrics by Template

[See](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-binding-group/deliverability-metrics-by-template)

Provides aggregate metrics grouped by template over the time window specified.

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "template"``

### Deliverability Metrics by Watched Domain

[See](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-campaign/deliverability-metrics-by-watched-domain)

Provides aggregate metrics grouped by watched domain over the time window specified. The difference between domain and watched domain is that watched domains are comprised of the top 99% domains in the world.

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "watched-domain"``


### Time Series

[See](https://developers.sparkpost.com/api/#/reference/metrics/deliverability-metrics-by-template/time-series-metrics)

Provides deliverability metrics ordered by a precision of time.The following table describes the validation for the precision parameter:

* ``./sp-deliverability-metrics-cli --from "2014-02-01T00:00" --command "watched-domain"` --precision day`

| Value of  | Valid for time window of  |
|---|:-:|
|1min, 5min|	day|
|hour	| month|
|day, month | any|





### Deliverability Usage

```
--baseurl, -u            Optional baseUrl for SparkPost. [$SPARKPOST_BASEURL]
--apikey, -k 			 Required SparkPost API key [$SPARKPOST_API_KEY]
--username 				 Username this is a special case it is more common to use apikey
--password, -p 			 Username this is a special it is more common to use apikey
--verbose "false"		 Dumps additional information to console
--command "domain"		 Optional one of domain, binding, binding-group, campaign, template, watched-domain, time-series
--from, -f 				 Required Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T08:00. Default: One hour ago
--to 					 Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T00:00. Default: now.
--domains, -d 			 Optional Comma-delimited list of domains to include Example: gmail.com,yahoo.com,hotmail.com.
--campaigns, -c 		 Optional Comma-delimited list of campaigns to include. Example: Black Friday
--metrics, -m            Required Comma-delimited list of metric name for filtering
--templates 			 Optional comma-delimited list of template IDs to include Example: summer-sale
--nodes 				 Optional comma-delimited list of nodes to include ( Note: SparkPost Elite only ) Example: Email-MSys-1,Email-MSys-2,Email-MSys-3
--bindings 				 Optional comma-delimited list of bindings to include (Note: SparkPost Elite only) Example: Confirmation
--binding_groups 		 Optional comma-delimited list of binding groups to include (Note: SparkPost Elite only) Example: Transaction
--protocols 			 Optional comma-delimited list of protocols for filtering (Note: SparkPost Elite only) Example: smtp
--timezone 				 Standard timezone identification string, defaults to UTC Example: America/New_York.
--limit 				 Optional maximum number of results to return Example: 5
--order_by 				 Optional metric by which to order results Example: count_injected
--precision              Precision of timeseries data returned Example: day. Possible values:  1min , 5min , 15min , hour , 12hr , day , week , month .
--help, -h				 show help
--version, -v			 print the version

```

| Metric Name |
|---|
|count_injected|
|count_bounce|
|count_rejected|
|count_delivered|
|count_delivered_first|
|count_delivered_subsequent|
|total_delivery_time_first|
|total_delivery_time_subsequent|
|total_msg_volume|
|count_policy_rejection|
|count_generation_rejection|
|count_generation_failed|
|count_inband_bounce|
|count_outofband_bounce|
|count_soft_bounce|
|count_hard_bounce|
|count_block_bounce|
|count_admin_bounce|
|count_undetermined_bounce|
|count_delayed|
|count_delayed_first|
|count_rendered|
|count_unique_rendered|
|count_unique_confirmed_opened|
|count_clicked|
|count_unique_clicked|
|count_targeted|
|count_sent|
|count_accepted|
|count_spam_complaint|	


## Message Events


| Compatibility | Compatible? |
|---|:-:|
|SparkPost| Yes |
|SparkPost Elite| Yes |
|Momentum| No |

The following options are available for the Message Event CLI:

| Option | Default | Descrption|
|---|:-:|:-:|
|--baseurl, -u| "https://api.sparkpost.com"|	Optional baseUrl for SparkPost. [$SPARKPOST_BASEURL]|
|--apikey, -k | |Required SparkPost API key [$SPARKPOST_API_KEY]|
|--username| |Username it is more common to use apikey|
|--password, -p| |Username it is more common to use apikey|
|--verbose| "false"|Dumps additional information to console|
|--bounce_classes, -b| |Optional comma-delimited list of bounce classification codes to search.|
|--campaign_ids, -i| |Optional comma-delimited list of campaign ID's to search. Example: "Example Campaign Name"|
|--events, -e||Optional comma-delimited list of event types to search. Defaults to all event types.|
|--friendly_froms|||Optional comma-delimited list of friendly_froms to search|
|--from, -f||Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T08:00. Default: One hour ago|
|--message_ids||Optional Comma-delimited list of message ID's to search. Example: 0e0d94b7-9085-4e3c-ab30-e3f2cd9c273e.|
|--page||Optional results page number to return. Used with per_page for paging through result. Example: 25. Default: 1|
|--per_page||Optional number of results to return per page. Must be between 1 and 10,000 (inclusive). Example: 100. Default: 1000.|
|--reason||Optional bounce/failure/rejection reason that will be matched using a wildcard (e.g., %%reason%%). Example: bounce.|
|--recipients||Optional Comma-delimited list of recipients to search. Example: recipient@example.com|
|--template_ids||Optional Comma-delimited list of template ID's to search. Example: templ-1234.|
|--timezone||Optional Standard timezone identification string. Example: America/New_York. Default: UTC|
|--to||Optional Datetime in format of YYYY-MM-DDTHH:MM. Example: 2016-02-10T00:00. Default: now.|
|--transmission_ids||Optional Comma-delimited list of transmission ID's to search (i.e. id generated during creation of a transmission). Example: 65832150921904138.|
