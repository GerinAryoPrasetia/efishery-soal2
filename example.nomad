job "test-efishery-api-gerinaryo" {
  datacenters = ["dc1"]

  group "cache" {
    network {
      port "db" {
        to = 1323
      }
    }

    task "redis" {
      driver = "docker"

      config {
        image = "gerinaryo/test-efishery-go-api:latest"

        ports = ["db"]
      }

      resources {
        cpu    = 500
        memory = 256
      }
    }
  }
}
