consul {
  address = "localhost:8500"

  retry {
    enabled  = true
    attempts = 12
    backoff  = "250ms"
  }
}

template {
  source      = "./nginx/templates/load-balancer.conf.ctmpl"
  destination = "./nginx/nginx.conf"
  perms       = 0644
  command     = "docker exec nginx-zookeeper-nginx-1 nginx -s reload"
}
