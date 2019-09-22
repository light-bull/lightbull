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

### Create show

    curl -X POST -d '{"name":"Show Name"}' 'http://localhost:8080/api/shows'

### Get show details

    curl -X GET 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Update show

    curl -X POST -d '{"name":"New Show Name"}' 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Delete show

    curl -X DELETE 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

## Visuals

### Create visual

    curl -X POST -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1", "name":"Visual Name"}' 'http://localhost:8080/api/visuals'

### Get visual details

    curl -X GET 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Update visual

TODO

### Delete visual

    curl -X DELETE 'http://localhost:8080/api/shows/61370850-aa63-44f7-a9d9-49b6292763b8'

### Get all visual names (of all shows)

    curl -X GET 'http://localhost:8080/api/visuals'

## Parameters

TODO: GET and POST for values
