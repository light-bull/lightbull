# Authentication

    curl -X POST -d '{"password":"lightbull"}' 'http://localhost:8080/api/auth'

Since the JWT is necessary for all following requests, it is recommended to write it to an environment variable for development and debugging:

    export jwt=$(curl -X POST -d '{"password":"lightbull"}' 'http://localhost:8080/api/auth' | jq -r '.jwt')

# Config

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/config'

# System

## Shutdown

    curl -H "Authorization: Bearer ${jwt}" -X POST http://localhost:8080/api/shutdown

## Ethernet configuration

### Get configuration

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/ethernet'

### Change configuration

    curl -H "Authorization: Bearer ${jwt}" -X POST -d '{"mode":"static","ip":"10.0.0.10/24","gateway":"10.0.0.1","dns":"8.8.8.8"}' 'http://localhost:8080/api/ethernet'

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

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/shows'

### Create show

    curl -H "Authorization: Bearer ${jwt}" -X POST -d '{"name":"Show Name"}' 'http://localhost:8080/api/shows'

### Get show details

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Update show

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"name":"New Show Name", "favorite": true}' 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

### Delete show

    curl -H "Authorization: Bearer ${jwt}" -X DELETE 'http://localhost:8080/api/shows/4f7f6045-bd3f-4fa3-9790-008df78571c1'

## Visuals

### Create visual

    curl -H "Authorization: Bearer ${jwt}" -X POST -d '{"showId":"4f7f6045-bd3f-4fa3-9790-008df78571c1", "name":"Visual Name"}' 'http://localhost:8080/api/visuals'

### Get visual details

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Update visual

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"name":"New Visual Name"}' 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Delete visual

    curl -H "Authorization: Bearer ${jwt}" -X DELETE 'http://localhost:8080/api/visuals/61370850-aa63-44f7-a9d9-49b6292763b8'

### Get all visual names (of all shows)

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/visuals'

## Groups

### Add group

    curl -H "Authorization: Bearer ${jwt}" -X POST -d '{"visualId":"61370850-aa63-44f7-a9d9-49b6292763b8", "parts":["horn_left", "horn_right"], "effectType":"singlecolor"}' 'http://localhost:8080/api/groups'

### Update group

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"parts": ["horn_left"], "effectType":"othereffect"}' 'http://localhost:8080/api/groups/e8a6b7c4-d2fe-4701-9d73-fe2e8377d0fb'

It's possible to set only "parts" or "effect".

### Delete group

    curl -H "Authorization: Bearer ${jwt}" -X DELETE 'http://localhost:8080/api/groups/e8a6b7c4-d2fe-4701-9d73-fe2e8377d0fb'

## Parameters

### Get parameter

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf'

### Update parameter

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"current":{"r":255, "g":0, "b":0}}' 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf' 
    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"default":{"r":255, "g":0, "b":0}}' 'http://localhost:8080/api/parameters/53d84761-d08f-4ef5-8ec2-5692d9a1a8cf' 

The current and default value can also be set in the same request.

### Restore default

    TODO

## Current show and visual

### Get current show and visual

    curl -H "Authorization: Bearer ${jwt}" -X GET 'http://localhost:8080/api/current'

### Set current show or visual

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1"}' 'http://localhost:8080/api/current'
    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"visual":"61370850-aa63-44f7-a9d9-49b6292763b8"}' 'http://localhost:8080/api/current'
    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"show":"4f7f6045-bd3f-4fa3-9790-008df78571c1","visual":"61370850-aa63-44f7-a9d9-49b6292763b8"}' 'http://localhost:8080/api/current'

When the show is changed and no visual is specified, the current visual is set to null (but only if the value changes). If only a visual is specified, it must belong to the current show.

### Set blank

    curl -H "Authorization: Bearer ${jwt}" -X PUT -d '{"blank":"true"}' 'http://localhost:8080/api/current'

Sets the visual to null which means that the LEDs are off. The current show is not changed.

# Websockets

## Connect

Connect to `ws://localhost:8080/api/ws`. The first message has to be:

    {"topic":"identify","payload":{"token":"$jwt"}}

From there on, the client will receive updates. The returned connection ID should be included in all future HTTP API requests in the `X-Lightbull-Connection-Id` header.

## Update parameters

    {"topic":"parameter","payload":{"id":"a5922724-f395-4a43-b38c-8b78de0ec2be","value":{"r": 128,"g": 255,"b": 255}}}current_show