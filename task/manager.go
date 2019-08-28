package task

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"
)

type Status string

const (
	DoesNotExist           Status = "does_not_exist"
	PendingRegistration    Status = "pending_registration"
	WaitForAPIKey          Status = "wait_for_api_key"
	Registered             Status = "registered"
	WaitForTaskDescription Status = "wait_for_task_desc"

	dumpFileName = "dump.json"
)

type user struct {
	Username string
	APIKey   string
	Status   Status
}

func (u *user) getStatus() Status {
	return u.Status
}

func (u *user) setStatus(s Status) {
	u.Status = s
}

func (u *user) getAPIKey() string {
	return u.APIKey
}

type Manager struct {
	users map[string]*user
	lock  sync.RWMutex
}

func NewManager() (*Manager, error) {
	users := make(map[string]*user)

	b, err := ioutil.ReadFile(dumpFileName)
	if err != nil {
		log.Printf("err: %v", err)

		return &Manager{
			users: users,
		}, nil
	}

	err = json.Unmarshal(b, &users)
	if err != nil {
		log.Printf("err: %v", err)

		return &Manager{
			users: users,
		}, nil
	}

	return &Manager{
		users: users,
	}, nil
}

func (m *Manager) Dump() {
	m.lock.RLock()
	b, err := json.Marshal(m.users)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
	m.lock.RUnlock()

	err = ioutil.WriteFile(dumpFileName, b, 0644)
	if err != nil {
		log.Printf("err: %v", err)
		return
	}
}

func (m *Manager) Add(username string) {
	m.lock.Lock()
	m.users[username] = &user{
		Username: username,
		APIKey:   "",
		Status:   PendingRegistration,
	}
	m.lock.Unlock()
}

func (m *Manager) Register(username, apikey string) {
	m.lock.Lock()
	m.users[username] = &user{
		Username: username,
		APIKey:   apikey,
		Status:   Registered,
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
}

func (m *Manager) GetAPIKey(username string) string {
	m.lock.RLock()
	defer m.lock.RUnlock()

	user, ok := m.users[username]
	if !ok {
		// should not reach here
		log.Println("err: user not found while updating status")
		return ""
	}

	return user.getAPIKey()
}
