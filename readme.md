# Logsend

[![Build Status](https://travis-ci.org/ezotrank/logsend.svg?branch=master)](https://travis-ci.org/ezotrank/logsend)

---
Logsend is high-performance tool for processing logs.


## What is it

This like [Logstash](http://logstash.net) but more tiny and written by [Golang](http://golang.org).
Supported outputs:

* [Influxdb](#influxdb)
* [Statsd](#statsd)
* [MySQL](#mysql)


##<a name="instalation"></a>Instalation

Just download binary by your platform from [GitHub Latest Releases](https://github.com/ezotrank/logsend/releases/latest) and unzip, or use this installer

Install logsend to `/usr/local/bin`

```
curl -L http://logsend.io/get|sudo bash
```
OR to your own directory

```
curl -L http://logsend.io/get|bash -s /tmp
```

##<a name="how_use"></a>How it can be used

As daemon, watching file in directory:

```
logsend -config=config.json /logs
```

Using PIPE:

```
tail -F /logs/*.log |logsend -config=config.json
```

```
cat /logs/*.log |logsend -config=config.json
```

```
ssh user@host "tail -F /logs/*.log"|logsend -config=config.json
```

Or using PIPE without config.json:

```
tail -F /logs/*.log |logsend -influx-dbname test -influx-host 'hosta:4444' -regex='\d+'
```

Daemonize:

```
logsend -config=config.json /logs 2> error.log &
```

## Benchmarks
<a name="benchmarks"></a>

### With 1 rule

| Log files | Lines per log | Result (real/user/sys)       |
|:----------:|:-------------:|:----------------------------:|
| 1          | 250k          | 0m3.186s 0m3.742s 0m1.672s   |
| 5          | 250k          | 0m5.641s 0m13.423s 0m1.015s  |
| 10         | 250k          | 0m8.577s 0m25.781s 0m1.341s  |

***config.json***

```
{
  "influxdb": {
    "host": "localhost:4444",
    "user": "root",
    "password": "root",
    "database": "logers",
    "udp": true,
    "send_buffer": 8
  },
  "groups": [
    {
        "mask": "test.log",
        "rules": [
            {
                "regexp": "test string (?P<word_STRING>\\\\w+)",
                "influxdb": {
                    "name": "test"
                }
            }
        ]
    }
  ]
}
```

### With 3 rules

| Log files | Lines per log | Result (real/user/sys)       |
|:----------:|:-------------:|:----------------------------:|
| 1          | 250k          | 0m5.722s 0m7.638s 0m3.779s   |
| 5          | 250k          | 0m9.676s 0m30.553s 0m2.410s  |
| 10         | 250k          | 0m18.385s 0m58.795s 0m2.534s |


***config.json***

```
{
  "influxdb": {
    "host": "localhost:4444",
    "user": "root",
    "password": "root",
    "database": "logers",
    "udp": true,
    "send_buffer": 8
  },
  "groups": [
    {
        "mask": "test.log",
        "rules": [
            {
                "regexp": "test string (?P<word_STRING>\\\\w+)",
                "influxdb": {
                    "name": "test"
                }
            },
            {
                "regexp": "(?P<word_STRING>\\\\w+) string one",
                "influxdb": {
                    "name": "test"
                }
            },
            {
                "regexp": "string (?P<word_STRING>\\\\w+) one",
                "influxdb": {
                    "name": "test"
                }
            }
        ]
    }
  ]
}
```


## Starting

```
logsend -config config.json /logs
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

## Outputs

###<a name="influxdb"></a>Influxdb

config.json:

```
{
  "influxdb": {
    "host": "influxdbhost:8086",
    "user": "root",
    "password": "root",
    "database": "logers",
    "send_buffer": 12
  },
  "groups": [
    {
        "mask": "some.[0-9]+.log$",
        "rules": [
            {
                "regexp": "\\[W .+ chain_processor.+\\] \\[.+\\] (?P<gate_STRING>\\w+_gate) -> (?P<exec_time_FLOAT>\\d+.\\d+)",
                "influxdb": {
                    "name": "gate_response",
                    "extra_fields": [
                      ["host", "HOST"]
                    ]
                }
            },
            {
                "regexp": "^\\[I .+ web.+\\] (?P<code_INT>\\d+) \\w+ (?P<name_STRING>\/\\S+)(?:\\?.+) \\(.+\\) (?P<resp_time_FLOAT>\\d+.\\d+)ms",
                "influxdb": {
                    "name": "adaptors",
                    "extra_fields": [
                      ["host", "HOST"]
                    ]
                }
            }
        ]
    }
  ]
}
```

another examples how to run:

only one rule:

```
tail -F some.log| logsend -influxdb-host "influxdbhost:8086"
-influxdb-user root -influxdb-password root
-influxdb-database logers -influxdb-send_buffer 12
-influxdb-extra_fields 'host,HOST'
-regex '\[W .+ chain_processor.+\] \[.+\] (?P<gate_STRING>\w+_gate) -> (?P<exec_time_FLOAT>\d+.\d+)'
```

using previous config file:

```
tail -F some.log| logsend -config config.json
```

###<a name="statsd"></a>Statsd

config.json:

```
{
  "statsd": {
    "host": "statsdhost:8125",
    "prefix": "test",
    "interval": "1s"
  },
  "groups": [
    {
        "mask": "some.[0-9]+.log$",
        "rules": [
            {
                "regexp": "\\[W .+ chain_processor.+\\] \\[.+\\] (?P<gate_STRING>\\w+_gate) -> (?P<exec_time_FLOAT>\\d+.\\d+)",
                "statsd": {
                  "timing": [
                    ["gate_exec_time", "exec_time"]
                  ]
                }
            },
            {
                "regexp": "^\\[I .+ web.+\\] (?P<code_INT>\\d+) \\w+ (?P<name_STRING>\/\\S+)(?:\\?.+) \\(.+\\) (?P<resp_time_FLOAT>\\d+.\\d+)ms",
                "statsd": {
                  "increment": [
                  	"adaptors_exec_count"
                  ]
                }
            }
        ]
    }
  ]
}
```

another examples how to run:

only one rule:

```
tail -F some.log| logsend -statsd-host "statsdhost:8125"
-statsd-host "statsdhost:8125" -statsd-prefix test
-statsd-interval 1s -statsd-timing 'gate_exec_time,exec_time'
-regex '\[W .+ chain_processor.+\] \[.+\] (?P<gate_STRING>\w+_gate) -> (?P<exec_time_FLOAT>\d+.\d+)'
```

using previous config file:

```
tail -F some.log| logsend -config config.json
```

###<a name="mysql"></a>MySQL

config.json:

```
{
  "mysql": {
    "host": "root:toor@/test1?timeout=30s&strict=true"
  },
  "groups": [
    {
        "mask": "some.[0-9]+.log$",
        "rules": [
            {
                "regexp": "\\[W .+ chain_processor.+\\] \\[.+\\] (?P<gate_STRING>\\w+_gate) -> (?P<exec_time_FLOAT>\\d+.\\d+)",
                "mysql": {
                  "query": [
                    "insert into test1(teststring, testfloat) values('{{.gate}}', {{.exec_time}});",
                    "insert into test2(teststring, testfloat) values('{{.gate}}', {{.exec_time}});"
                  ]
                }
            },
            {
                "regexp": "^\\[I .+ web.+\\] (?P<code_INT>\\d+) \\w+ (?P<name_STRING>\/\\S+)(?:\\?.+) \\(.+\\) (?P<resp_time_FLOAT>\\d+.\\d+)ms",
                "mysql": {
                  "query": [
                    "insert into test1(teststring, testfloat) values('a', 22);"
                  ]
                }
            }
        ]
    }
  ]
}
```

another examples how to run:

only one rule:

```
tail -F some.log| logsend -mysql-host "root:toor@/test1?timeout=30s&strict=true"
-mysql-query "insert into test1(teststring, testfloat) values('{{.gate}}', {{.exec_time}});"
-regex '\[W .+ chain_processor.+\] \[.+\] (?P<gate_STRING>\w+_gate) -> (?P<exec_time_FLOAT>\d+.\d+)'
```

using previous config file:

```
tail -F some.log| logsend -config config.json
```

## Tips

* use flag `-debug` for more info
* use flag `-dry-run` for processing logs but not send to destination
* use flag `-read-whole-log` for reading whole log file and continue reading
* use flag `-read-once` better for use with -read-whole-log, just read whole log and exit