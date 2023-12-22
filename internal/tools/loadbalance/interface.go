package loadbalance

type LoadBalance interface {
	Add(...string) error
	Get(string) (string, error)

	// Update use for service-discovery to update the nodes
	Update()
}
