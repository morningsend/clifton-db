/*
Package cluster manager is configured to be a Raft learner in all multi-raft groups
Exactly one cluster manager instance is deployed to a CliftonDB cluster, and it only
replicates configuration change messages.

A cluster manager provides partition discover mechanism which is used by routers
in all nodes to receive most up-to-date routing information.

A cluster manager implements its own storage.

 */
package cluster_manager
