package callgraphutil

import "golang.org/x/tools/go/callgraph"

type NodeID int

type discoveredNode struct {
	id      NodeID
	callers []NodeID
	callees []NodeID
	visited bool
	marked  bool
}

func CollectTargetNodesBySourceNode(sourceNodes map[NodeID]callgraph.Node, targetNodes map[NodeID]callgraph.Node) map[NodeID][]NodeID {

	discoveredNodes := make(map[NodeID]*discoveredNode)
	targetNodesByNode := make(map[NodeID][]NodeID)

	{
		remainingNodes := make(map[NodeID]callgraph.Node)
		for nodeID, node := range targetNodes {
			discoveredNodes[nodeID] = &discoveredNode{
				id:      nodeID,
				callers: []NodeID{},
				callees: []NodeID{},
			}
			remainingNodes[NodeID(node.ID)] = node
		}
		for len(remainingNodes) > 0 {
			for nodeID, node := range remainingNodes {
				if callee, _ := discoveredNodes[nodeID]; !callee.visited {
					for _, call := range node.In {
						if call.Caller.ID != node.ID /* omit recursive calls */ {
							callerID := NodeID(call.Caller.ID)
							caller, found := discoveredNodes[callerID]
							if !found {
								caller = &discoveredNode{id: callerID, callers: []NodeID{}, callees: []NodeID{}}
								discoveredNodes[callerID] = caller
							}
							callee.callers = append(callee.callers, caller.id)
							caller.callees = append(caller.callees, nodeID)
							if _, found := sourceNodes[callerID]; !found {
								remainingNodes[callerID] = *call.Caller
							}
						}
					}
					callee.visited = true
				}
				delete(remainingNodes, nodeID)
				break
			}
		}
	}
	{
		remainingNodes := make(map[NodeID]bool)
		for nodeID := range sourceNodes {
			if _, ok := discoveredNodes[nodeID]; ok {
				remainingNodes[nodeID] = true
			}
		}
		for len(remainingNodes) > 0 {
			for nodeID := range remainingNodes {
				if node := discoveredNodes[nodeID]; !node.marked {
					node.marked = true
					for _, calleeID := range node.callees {
						remainingNodes[calleeID] = true
					}
				}
				delete(remainingNodes, nodeID)
				break
			}
		}
	}
	{
		remainingNodes := make(map[NodeID]bool)
		for nodeID := range targetNodes {
			if node := discoveredNodes[nodeID]; node.marked {
				targetNodesByNode[nodeID] = []NodeID{nodeID}
				remainingNodes[nodeID] = true
			}
		}
		for len(remainingNodes) > 0 {
			for nodeID := range remainingNodes {
				node := discoveredNodes[nodeID]
				for _, callerID := range node.callers {
					if caller := discoveredNodes[callerID]; caller.marked {
						existingTargets := targetNodesByNode[callerID]
						additionalTargets := make([]NodeID, 0)
					targets:
						for _, targetID := range targetNodesByNode[nodeID] {
							for _, existingTargetID := range existingTargets {
								if targetID == existingTargetID {
									continue targets
								}
							}
							additionalTargets = append(additionalTargets, targetID)
						}
						if len(additionalTargets) > 0 {
							targetNodesByNode[callerID] = append(existingTargets, additionalTargets...)
							remainingNodes[callerID] = true
						}
					}
				}
				delete(remainingNodes, nodeID)
				break
			}
		}
	}

	return targetNodesByNode
}
