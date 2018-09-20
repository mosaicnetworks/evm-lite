package config

import (
	"fmt"
	"time"

	_raft "github.com/hashicorp/raft"
)

var (
	defaultRaftDir     = fmt.Sprintf("%s/raft", DefaultDataDir)
	defaultSnapshotDir = fmt.Sprintf("%s/snapshots", defaultRaftDir)
	defaultRaftID      = defaultNodeAddr
)

//RaftConfig contains the configuration of a Raft node
type RaftConfig struct {

	// ProtocolVersion allows a Raft server to inter-operate with older
	// Raft servers running an older version of the code. This is used to
	// version the wire protocol as well as Raft-specific log entries that
	// the server uses when _speaking_ to other servers. There is currently
	// no auto-negotiation of versions so all servers must be manually
	// configured with compatible versions. See ProtocolVersionMin and
	// ProtocolVersionMax for the versions of the protocol that this server
	// can _understand_.
	ProtocolVersion _raft.ProtocolVersion `mapstructure:"protocol_version"`

	// HeartbeatTimeout specifies the time in follower state without
	// a leader before we attempt an election.
	HeartbeatTimeout time.Duration `mapstructure:"heartbeat"`

	// ElectionTimeout specifies the time in candidate state without
	// a leader before we attempt an election.
	ElectionTimeout time.Duration `mapstructure:"election_timeout"`

	// CommitTimeout controls the time without an Apply() operation
	// before we heartbeat to ensure a timely commit. Due to random
	// staggering, may be delayed as much as 2x this value.
	CommitTimeout time.Duration `mapstructure:"commit_timeout"`

	// MaxAppendEntries controls the maximum number of append entries
	// to send at once. We want to strike a balance between efficiency
	// and avoiding waste if the follower is going to reject because of
	// an inconsistent log.
	MaxAppendEntries int `mapstructure:"max_append_entries"`

	// If we are a member of a cluster, and RemovePeer is invoked for the
	// local node, then we forget all peers and transition into the follower state.
	// If ShutdownOnRemove is is set, we additional shutdown Raft. Otherwise,
	// we can become a leader of a cluster containing only this node.
	ShutdownOnRemove bool `mapstructure:"shutdown_on_remove"`

	// TrailingLogs controls how many logs we leave after a snapshot. This is
	// used so that we can quickly replay logs on a follower instead of being
	// forced to send an entire snapshot.
	TrailingLogs uint64 `mapstructure:"trailing_logs"`

	// SnapshotInterval controls how often we check if we should perform a snapshot.
	// We randomly stagger between this value and 2x this value to avoid the entire
	// cluster from performing a snapshot at once.
	SnapshotInterval time.Duration `mapstructure:"snapshot_interval"`

	// SnapshotThreshold controls how many outstanding logs there must be before
	// we perform a snapshot. This is to prevent excessive snapshots when we can
	// just replay a small set of logs.
	SnapshotThreshold uint64 `mapstructure:"snapshot_threshold"`

	// LeaderLeaseTimeout is used to control how long the "lease" lasts
	// for being the leader without being able to contact a quorum
	// of nodes. If we reach this interval without contact, we will
	// step down as leader.
	LeaderLeaseTimeout time.Duration `mapstructure:"leader_lease_timeout"`

	// StartAsLeader forces Raft to start in the leader state. This should
	// never be used except for testing purposes, as it can cause a split-brain.
	StartAsLeader bool `mapstructure:"start_as_leader"`

	// The unique ID for this server across all time. When running with
	// ProtocolVersion < 3, you must set this to be the same as the network
	// address of your transport.
	LocalID _raft.ServerID `mapstructure:"server_id"`

	/*------------------------------------------------------------------------*/

	//XXX TODO improve this

	RaftDir     string `mapstructure:"dir"`
	SnapshotDir string `mapstructure:"snapshot_dir"`
	NodeAddr    string `mapstructure:"node_addr"`
}

//DefaultRaftConfig returns the default configuration for a Raft node
func DefaultRaftConfig() *RaftConfig {
	return &RaftConfig{
		ProtocolVersion:    _raft.ProtocolVersionMax,
		HeartbeatTimeout:   1000 * time.Millisecond,
		ElectionTimeout:    1000 * time.Millisecond,
		CommitTimeout:      50 * time.Millisecond,
		MaxAppendEntries:   64,
		ShutdownOnRemove:   true,
		TrailingLogs:       10240,
		SnapshotInterval:   120 * time.Second,
		SnapshotThreshold:  8192,
		LeaderLeaseTimeout: 500 * time.Millisecond,
		LocalID:            _raft.ServerID(defaultRaftID),
		RaftDir:            defaultRaftDir,
		SnapshotDir:        defaultSnapshotDir,
		NodeAddr:           defaultNodeAddr,
	}
}
