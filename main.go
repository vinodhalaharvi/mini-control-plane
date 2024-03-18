package main

func main() {
	server := NewKubeAPIServer()
	controller := NewKubeControllerManager()

	go controller.Run()
	server.Run()
}
