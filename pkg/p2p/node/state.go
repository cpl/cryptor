package node

// StateNode defines an alias of the underlying type used for enumerating
// the possible node states.
type StateNode byte

const (
	// StateStopped is the default starting state of a Node, or the state
	// reached after calling node.Stop(). During this state node operations
	// are limited. A node also can't be connected during this state.
	StateStopped StateNode = iota

	// StateRunning is the state after successfully calling node.Start() or
	// node.Disconnect(). During this state the node allows for some
	// configurations or actions to be taken.
	StateRunning

	// StateConnected is entered after calling node.Connect() on a running node.
	// Node must be running in order to also be connected to the network.
	StateConnected
)

// a simple lookup for state "ID" to state string, allows for nice displaying
var stateNames = map[StateNode]string{
	StateStopped:   "STOPPED",
	StateRunning:   "RUNNING",
	StateConnected: "CONNECTED",
}

// String will return the name of the current state as a string from stateNames.
func (s StateNode) String() string {
	return stateNames[s]
}

// other possible future states: Crashed, Restarting, Blocked, Idle

// State returns a node state/status. This state can be one of the
// following constant pre-defined byte values:
// - Stopped
// - Running
// - Connected
// Other values may be added in the future.
func (n *Node) State() StateNode {
	// lock states for reading
	n.state.RLock()
	defer n.state.RUnlock()

	// check state and return
	switch {
	case n.state.isConnected:
		return StateConnected
	case n.state.isRunning:
		return StateRunning
	default:
		return StateStopped
	}
}
