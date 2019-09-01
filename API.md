# System

## Shutdown

Supported methods: POST

    curl -X POST http://localhost:8080/api/shutdown

## Ethernet configuration

Supported methods: GET, POST

    curl -X POST -d '{"mode":"static","ip":"10.0.0.10/24","gateway":"10.0.0.1","dns":"8.8.8.8"}' 'http://localhost:8080/api/ethernet'

### Details
Key     | Description
--------|---------------------
mode    | Mode: down, static, dhcp-client, dhcp-server, unmanaged (unmanaged cannot be set)
ip      | IP and subnet in CIDR notation ("192.168.0.1/24", required for static and dhcp-server)
gateway | IP address of gateway
dns     | IP address of DNS server

# Shows

## Shows

Supported methods: GET, POST, PUT

### Get all shows

    curl -X GET 'http://localhost:8080/api/shows'

### Get show details

    curl -X GET 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Create show

    curl -X POST -d '{"name":"Show Name"}' 'http://localhost:8080/api/shows

### Update show

    curl -X PUT -d '{"name":"New Show Name", "id":"4f7f6045-bd3f-4fa3-9790-008df78571c1"}' 'http://localhost:8080/api/shows'

## Visuals

TODO: GET should return all related data (visual meta data, groups, effects, parameters)

## Parameters

TODO: GET and POST for values
