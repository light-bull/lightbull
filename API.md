# Authentication

    curl -X POST -d '{"password":"lightbull"}' 'http://localhost:8080/api/auth'

# Config

    curl -X GET 'http://localhost:8080/api/config'

# System

## Shutdown

    curl -X POST http://localhost:8080/api/shutdown

## Ethernet configuration

### Get configuration

    curl -X GET 'http://localhost:8080/api/ethernet'

### Change configuration

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

### Get all shows

    curl -X GET 'http://localhost:8080/api/shows'

### Create show

    curl -X POST -d '{"name":"Show Name"}' 'http://localhost:8080/api/shows'

### Get show details

    curl -X GET 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Update show

    curl -X PUT -d '{"name":"New Show Name", "favorite": true}' 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Delete show

    curl -X DELETE 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

## Visuals

### Create visual

    curl -X POST -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1", "name":"Visual Name"}' 'http://localhost:8080/api/visuals'

### Get visual details

    curl -X GET 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Update visual

    curl -X PUT -d '{"name":"New Visual Name"}' 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Delete visual

    curl -X DELETE 'http://localhost:8080/api/shows/61370850-aa63-44f7-a9d9-49b6292763b8'

### Get all visual names (of all shows)

    curl -X GET 'http://localhost:8080/api/visuals'

## Groups

### Add group

    curl -X POST -d '{"visual":"61370850-aa63-44f7-a9d9-49b6292763b8", "parts":["horn_left", "horn_right"], "effect":"singlecolor"}' 'http://localhost:8080/api/groups'

### Update group

    curl -X PUT -d '{"parts": ["horn_left"], "effect":"othereffect"}' 'http://localhost:8080/api/groups/e8a6b7c4-d2fe-4701-9d73-fe2e8377d0fb'

    It's possible to set only "parts" or "effect".

### Delete group

    curl -X DELETE 'http://localhost:8080/api/groups/e8a6b7c4-d2fe-4701-9d73-fe2e8377d0fb'

## Parameters

### Get parameter

    curl -X GET 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf'

### Update parameter

    curl -X PUT -d '{"current":{"r":255, "g":0, "b":0}}' 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf' 
    curl -X PUT -d '{"default":{"r":255, "g":0, "b":0}}' 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf' 

    The current and default value can also be set in the same request.

### Restore default

    TODO

## Current show and visual

### Get current show and visual

    curl -X GET 'http://localhost:8080/api/current'

### Set current show or visual

    curl -X PUT -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1"}' 'http://localhost:8080/api/current'
    curl -X PUT -d '{"visual":"61370850-aa63-44f7-a9d9-49b6292763b8"}' 'http://localhost:8080/api/current'
    curl -X PUT -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1","visual":"61370850-aa63-44f7-a9d9-49b6292763b8"}' 'http://localhost:8080/api/current'

    When the show is changed and no visual is specified, the current visual is set to null (but only if the value changes). If only a visual is specified, it must belong to the current show.

### Set blank

    curl -X PUT -d '{"blank":"true"}' 'http://localhost:8080/api/current'

    Sets the visual to null which means that the LEDs are off. The current show is not changed.