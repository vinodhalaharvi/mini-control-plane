package main

func main() {
	controller := NewKubeControllerManager()
	go controller.Run()

	server := NewKubeAPIServer()
	server.setupRoutes()
	server.Run()
}
