package application

type Application interface {
	// Setup setups the application
	Setup()
	// Run runs the application
	Run()
}
