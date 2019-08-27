package task

import (
	"log"
	"sync"
)

type Status string

const (
	DoesNotExist        Status = "does_not_exist"
	PendingRegistration Status = "pending_registration"
	WaitForAPIKey       Status = "wait_for_api_key"
	Registered          Status = "registered"
)

type user struct {
	username string
	apiKey   string
	status   Status
}

func (u *user) getStatus() Status {
	return u.status
}

func (u *user) setStatus(s Status) {
	u.status = s
}

type Manager struct {
	users map[string]*user
	lock  sync.RWMutex
}

func NewManager() (*Manager, error) {
	// TODO: load from db.
	return &Manager{
		users: make(map[string]*user),
	}, nil
}

func (m *Manager) Add(username string) {
	m.lock.Lock()
	m.users[username] = &user{
		username: username,
		apiKey:   "",
		status:   PendingRegistration,
	}
	m.lock.Unlock()
}

func (m *Manager) Register(username, apikey string) {
	m.lock.Lock()
	m.users[username] = &user{
		username: username,
		apiKey:   apikey,
		status:   Registered,
	}
	m.lock.Unlock()
}

func (m *Manager) GetStatus(username string) Status {
	m.lock.RLock()
	defer m.lock.RUnlock()

	user, ok := m.users[username]
	if !ok {
		return DoesNotExist
	}

	log.Printf("user get: %#v", user)

	return user.getStatus()
}

func (m *Manager) SetStatus(username string, status Status) {
	m.lock.RLock()
	defer m.lock.RUnlock()

	user, ok := m.users[username]
	if !ok {
		// should not reach here
		log.Println("err: user not found while updating status")
		return
	}

	user.setStatus(status)

	log.Printf("user set: %#v", user)
}
