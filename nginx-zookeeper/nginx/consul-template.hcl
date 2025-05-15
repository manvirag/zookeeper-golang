consul {
  address = "consul:8500"

  retry {
    enabled  = true
    attempts = 12
    backoff  = "250ms"
  }
}

template {
  source      = "/etc/consul-template/templates/load-balancer.conf.ctmpl"
  destination = "/etc/nginx/nginx.conf"
  perms       = 0644
  command     = "docker exec nginx-zookeeper-nginx-1 nginx -s reload"
}
