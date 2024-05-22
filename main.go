package main

import "madaurus/dev/material/app/server"

func main() {

	server.NewGracefulServer().Run().Wait()

}
