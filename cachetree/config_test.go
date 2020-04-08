package cachetree_test

var testConfig = `
{
		"Driver":"tree",
		"Marshaler": "json",
		"Config":{
			"Debug":true,
			"Alias":{
				"alias":"te"
			},
			"Root":{
				"Driver":"dummycache",
				"Marshaler": "json"
			},
			"Children":{
				"te":{
					"Driver":"dummycache",
					"Marshaler": "json",
					"TTL":4800
				},
				"test/test2":{
					"Driver":"dummycache",
					"Marshaler": "json",
					"TTL":2400
				},
				"test":{
					"Driver":"dummycache",
					"Marshaler": "json",
					"TTL":3600
				}
			}
		}
}`
