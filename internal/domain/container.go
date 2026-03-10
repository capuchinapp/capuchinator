package domain

type ContainerState struct {
	Status ContainerStateStatus
	Health ContainerStateHealth
}

type ContainerStateStatus string

const (
	ContainerStateStatusCreated    ContainerStateStatus = "created"
	ContainerStateStatusRunning    ContainerStateStatus = "running"
	ContainerStateStatusPaused     ContainerStateStatus = "paused"
	ContainerStateStatusRestarting ContainerStateStatus = "restarting"
	ContainerStateStatusRemoving   ContainerStateStatus = "removing"
	ContainerStateStatusExited     ContainerStateStatus = "exited"
	ContainerStateStatusDead       ContainerStateStatus = "dead"
)

func (s ContainerStateStatus) String() string {
	return string(s)
}

type ContainerStateHealth string

const (
	ContainerStateHealthNoHealthcheck ContainerStateHealth = "none"      // Indicates there is no healthcheck
	ContainerStateHealthStarting      ContainerStateHealth = "starting"  // Starting indicates that the container is not yet ready
	ContainerStateHealthHealthy       ContainerStateHealth = "healthy"   // Healthy indicates that the container is running correctly
	ContainerStateHealthUnhealthy     ContainerStateHealth = "unhealthy" // Unhealthy indicates that the container has a problem
)

func (s ContainerStateHealth) String() string {
	return string(s)
}
