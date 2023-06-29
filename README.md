# Purpose
Tool to manage loadbalancer orinigs in CloudFlare LoadBalancer pools.
For example if you need enable and disable a bunch of origins in order to maintain them.

# Usage
Define environment variables:
```
CF_API_KEY=12312312313221 # API KEY FROM CloudFlare
CF_API_EMAIL=user@example.com # You CloudFlare account email
CF_API_ACCOUNTID=1231321312313212 # You account ID, could be found in CloudFlare WebUI.
CF_API_ZONEID=example.com #Zone you want to manage (Optional)
```

Examples
```
./cf -update -pool example-test -origin app-1 -state disable
```

# Build

```
go build
```

