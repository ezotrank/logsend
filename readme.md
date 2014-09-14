# Logsend
---
Logsend is a tool for managing your logs.


## What is it

This like [Logstash](http://logstash.net) but more tiny and written by [Golang](http://golang.org).
At the current time support [Influxdb](http://influxdb.com) and [Statsd](https://github.com/etsy/statsd/) outputs.

## Instalation

Just download binary by your platform from [GitHub Latest Releases](https://github.com/ezotrank/logsend/releases/latest) and unzip.

## Starting

```
./logsend -watch-dir=~/some_logs_folder
```

## Configuration

* [Influxdb config.json](#influxdb_config)
* [Influxdb advanced config.json](#influxdb_advanced_config)
* [Statsd config.json](#statsd_config)
* [Full config.json](#full_config)

**Influxdb configuration**

Description:

`"mask": "search.log$"` - match logs in watch-dir directory
`"regexp": "metric dbcount (?P<count_INT>[0-9]+)$",` - regexp for matching line. But we want to store some values from this line and we should use named group with type association, in this case `count_INT` will be perform to field `count` with type `int`. Also supported notation `FLOAT`, `STRING`.
`"name": "keys_count"` - database name in Influxdb

***<a name="influxdb_config"></a>config.json:***

```
{
    "influxdb": {
        "host": "host:4444",
        "user": "root",
        "password": "root",
        "database": "logers"
    },
    "groups": [
        {
            "mask": "search.log$",
            "rules": [
                {
                    "regexp": "metric dbcount (?P<count_INT>[0-9]+)$",
                    "influxdb": {
                        "name": "keys_count"
                    }
                }
            ]
        }
    ]
}
```

**Advanced Influxdb configuration**

Description:

`"udp": true,` - sending packets thought UDP(require that Influxdb supported this option)

`"send_buffer": 8` - messages per one packet. But don't forget, max size of packet should be no more then 1500 Byte

Some times we need added not captured params, and for this purpose we used `extra_fields` key.

```
"extra_fields": [
	["service", "search"],
	["host", "HOST"]
]
```

`["service", "search"],` - First is column name, second it's value.

`["host", "HOST"]` - HOST it's a special keyword who's returned hostname of machine.

***<a name="influxdb_advanced_config"></a>config.json***

```
{
    "influxdb": {
        "host": "host:4444",
        "user": "root",
        "password": "root",
        "database": "logers",
        "udp": true,
        "send_buffer": 8
    },
    "groups": [
        {
            "mask": "search.log$",
            "rules": [
                {
                    "regexp": "metric dbcount (?P<count_INT>[0-9]+)$",
                    "influxdb": {
                        "name": "keys_count",
                        "extra_fields": [
                        	["service", "search"],
                      		["host", "HOST"]
                    	]
                    }
                }
            ]
        }
    ]
}
```

**Statsd configuration**

Description:

`["search.prepare_proposal_time", "prepare_proposals"]` - set `prepare_proposals` to `search.prepare_proposal_time` as timing
`"search.prepare_proposals"` - increment by by this metric
`["search.keys_count", "prepare_proposals"]` - set non integer metrics 

***<a name="statsd_config"></a>config.json***

```
{
    "regexp": "PreapreProposals took (?P<prepare_proposals_DurationToMillisecond>.+)$",
    "statsd": {
        "timing": [
        	["search.prepare_proposal_time", "prepare_proposals"]
        ],
        "increment": [
            "search.prepare_proposals"
        ],
        "gauge": [
        	["search.keys_count", "prepare_proposals"]
        ]
    }
}
```

**Full Influx and Statsd configuration**

***<a name="full_config"></a>config.json***

```
{
  "influxdb": {
    "host": "hosta:4444",
    "user": "root",
    "password": "root",
    "database": "logers",
    "udp": true,
    "send_buffer": 8
  },
  "statsd": {
    "host": "hostb:8125",
    "prefix": "test.",
    "interval": "1s"
  },
  "groups": [
    {
        "mask": "search.log$",
        "rules": [
            {
                "regexp": "metric dbcount (?P<count_INT>[0-9]+)$",
                "influxdb": {
                    "name": "g_keys_count",
                    "extra_fields": [
                      ["service", "search"],
                      ["host", "HOST"]
                    ]
                },
                "statsd": {
                  "gauge": [
                    ["search.keys_count", "count"]
                  ]
                }
            },
            {
                "regexp": "PreapreProposals took (?P<prepare_proposals_DurationToMillisecond>.+)$",
                "influxdb": {
                    "name": "g_benchmark",
                    "extra_fields": [
                      ["service", "search"],
                      ["host", "HOST"]
                    ]
                },
                "statsd": {
                  "timing": [
                    ["search.prepare_proposal_time", "prepare_proposals"]
                  ],
                  "increment": [
                    "search.prepare_proposals"
                  ]
                }
            },
            {
                "regexp": "Completed (?P<code_INT>\\d+) .+ in (?P<tm_DurationToMillisecond>.+)$",
                "influxdb": {
                    "name": "g_responses",
                    "extra_fields": [
                      ["service", "search"],
                      ["host", "HOST"]
                    ]
                },
                "statsd": {
                  "increment": [
                    "search.requests"
                  ]
                }
            }
        ]
    },
    {
        "mask": "map.log$",
        "rules": [
            {
                "regexp": "metric dbcount (?P<count_INT>[0-9]+)$",
                "influxdb": {
                    "name": "g_keys_count",
                    "extra_fields": [
                      ["service", "map"],
                      ["host", "HOST"]
                    ]
                }
            },
            {
                "regexp": "PreapreProposals took (?P<prepare_proposals_DurationToMillisecond>.+)$",
                "influxdb": {
                    "name": "g_benchmark",
                    "extra_fields": [
                      ["service", "map"],
                      ["host", "HOST"]
                    ]
                }
            },
            {
                "regexp": "Completed (?P<code_INT>\\d+) .+ in (?P<tm_DurationToMillisecond>.+)$",
                "influxdb": {
                    "name": "g_responses",
                    "extra_fields": [
                      ["service", "map"],
                      ["host", "HOST"]
                    ]
                }
            }
        ]
    }
  ]
}
```


## Tips