package jobs

import (
	"fmt"
	"github.com/gocql/gocql"
	"github.com/gorhill/cronexpr"
	"github.com/prometheus/common/log"
)

// TODO * Add metrics for how often we try CAS because we did not see current
// TODO   correct job meta data.
// TODO * Load job info using EACH_QUROUM to make sure we pick up state correctly
// TODO   OTH, maybe a version info on the state will help with that without using
// TODO   EACH_QUORUM

// A repository implementation for Cassandra keyspace.
type CassandraRepo struct {
	session *gocql.Session
}

// getJobInfo retrieves the meta data for the job with the given name.
// Returns nil,nil for jobinfo if there is no job with the given name
// and nil,err if an error occurred.
func (m *CassandraRepo) GetJobinfo(name string) (*jobinfo, error) {
	li := jobinfo{name: name}
	var schedule string

	// Using CL LocalOne because the read is informational only - any data
	// we miss here will only cause an unnecessary locking attempt.
	err := m.session.
		Query("select enabled,lockttlsec,checksec,schedule,last,owner,state from jobs where name = ?", name).
		Consistency(gocql.LocalOne).
		Scan(&(li.enabled), &(li.lockttlsec), &(li.checksec), &schedule, &(li.last), &(li.owner), &(li.state))
	if err == gocql.ErrNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	li.schedule, err = cronexpr.Parse(schedule)
	if err != nil {
		return nil, err
	}

	return &li, nil
}

// saveState saves the given state for the given job.
// This does not check lock ownership.
func (m *CassandraRepo) SaveState(name string, state []byte) error {
	// Using CL EachQuorum to make sure that subsequent runs from any DC
	// pick up the correct state.
	if err := m.session.
		Query("update jobs set state = ?, last = unixTimestampOf(now()) where name = ?", state, name).
		Consistency(gocql.EachQuorum).
		Exec(); err != nil {
		// TODO map error
		return err
	}

	return nil
}

// touchLock updates the lock's TTL to ttlsec.
// This does not check ownership of the lock.
func (m *CassandraRepo) TouchLock(name string, owner string, ttlsec int) error {
	// Using CL LocalQuorum here because we want the TTL increase to be stored
	// safely in more than one node. The use of SERIAL in the lock CAS query will
	// make sure we see the TTL refresh when trying to get lock.
	return m.session.
		Query("update jobs using ttl ? set owner = ? where name = ?", ttlsec, owner, name).
		Consistency(gocql.LocalQuorum).
		Exec()
}

// tryGetLock attempts to acquire the given job's lock using a Cassandra CAS query.
// true will be returned if the lock has been acquired, false if it was already taken.
func (m *CassandraRepo) TryGetLock(name string, owner string, ttlsec int) (bool, error) {
	// CL of LocalOne is sufficient, because coordination happens based on CAS query.
	// Reads of older data in getJobInfo will only trigger unnecessary work.
	applied, err := m.session.
		Query("update jobs using ttl ? set owner = ? where name = ? if owner = null", ttlsec, owner, name).
		SerialConsistency(gocql.Serial).
		Consistency(gocql.LocalOne).
		ScanCAS()
	if err != nil {
		return false, err
	}
	return applied, nil
}

// Commit the given state to the job with the provided name
func (m *CassandraRepo) Commit(name string, state []byte) error {
	// Commit is an optimization, so LocalQuorum CL is enough - if data is lost
	// in worst case, it will only leat to re-doing some work.
	return m.session.
		Query("update locks set state = ? where name = ?", state, name).
		Consistency(gocql.LocalOne).
		Exec()
}

// log records a message in a per-job log messages table.
func (m *CassandraRepo) Log(name string, id string, event string, format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	if err := m.session.
		Query("update logs using ttl 864000 set id = ?,event = ?, msg = ? where name = ? and ts = now()",
			id, event, msg, name).
		Consistency(gocql.LocalOne).
		Exec(); err != nil {
		log.Errorf("Unable to save log (%s) '%s' to DB for %s,%s: %v)", event, msg, name, id, err)
	}
}

// CreateJob creates an entry in the jobs table for the provided jobinfo.
// This uses a CAS query to make sure it is only inserted once.
func (m *CassandraRepo) CreateJob(name string, jobcfg *JobCfg) error {
	_, err := m.session.
		Query("insert into jobs (name,enabled,lockttlsec,checksec,schedule) values (?,?,?,?,?) if not exists",
			name, jobcfg.Enabled, jobcfg.Lockttlsec, jobcfg.Checksec, jobcfg.Schedule).
		SerialConsistency(gocql.Serial).
		Consistency(gocql.EachQuorum).
		ScanCAS()
	return err
}
