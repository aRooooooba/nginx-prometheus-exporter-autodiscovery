# Based on [nginx-prometheus-exporter](https://github.com/nginxinc/nginx-prometheus-exporter)

It can auto discover the IPs of the tasks running on AWS ECS clusters.

## Build:
Binary:
```
go build -o nginx-prometheus-exporter-autodiscovery
```
Docker:
```
docker build -t nginx-prometheus-exporter-autodiscovery .
```

## Usage:
Binary:
```
./nginx-prometheus-exporter-autodiscovery
```
Or Docker:
```
docker run --rm -p 9113:9113 nginx-prometheus-exporter-autodiscovery
```
And then
```
curl localhost:9113/probe?cluster=<cluster name>
```
Acceptable url parameters:

`region`: mandatory, region of the AWS ECS cluster

`cluster`: mandatory, name of the AWS ECS cluster

`scheme`: optional, default `http`

`service`: optional, the name of the service where Nginx is running, default to cluster name

`port`: optional, default `443`

`path`: optional, the absolute path of the stub status endpoint, default `/stub_status`

Example:

```
curl 'localhost:9113/probe?region=us-west-2&cluster=my-cluster&port=8443&path=/basic_status'
```

## Permission

Required permission:

```
"ecs:ListTasks"
"ecs:DescribeTasks"
```
