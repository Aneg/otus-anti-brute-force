package ip_guard

import (
	"github.com/Aneg/otus-anti-brute-force/internal/constants"
	"github.com/Aneg/otus-anti-brute-force/internal/models"
	"github.com/Aneg/otus-anti-brute-force/internal/repositories"
	"net"
	"sync"
)

func NewMemoryIpGuard(listId constants.ListId, rep repositories.Masks) *MemoryIpGuard {
	return &MemoryIpGuard{
		listId: listId,
		rep:    rep,
		mu:     &sync.RWMutex{},
		masks:  make(map[int]*net.IPNet),
	}
}

type MemoryIpGuard struct {
	listId constants.ListId
	rep    repositories.Masks
	masks  map[int]*net.IPNet
	mu     *sync.RWMutex
}

func (m *MemoryIpGuard) AddMask(mask string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ipv4Net, err := net.ParseCIDR(mask)
	if err != nil {
		return false, err
	}

	for i := range m.masks {
		if m.masks[i].IP.String() == ipv4Net.IP.String() && m.masks[i].Mask.String() == ipv4Net.Mask.String() {
			// todo а оно работает?
			return false, nil
		}
	}

	model := models.Mask{
		Id:     0,
		Mask:   mask,
		ListId: m.listId,
	}

	if err = m.rep.Add(&model); err != nil {
		return false, err
	}
	m.masks[model.Id] = ipv4Net
	return true, nil
}

func (m *MemoryIpGuard) DropMask(mask string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, ipv4Net, err := net.ParseCIDR(mask)
	if err != nil {
		return false, err
	}

	for id := range m.masks {
		if m.masks[id].IP.String() != ipv4Net.IP.String() || m.masks[id].Mask.String() != ipv4Net.Mask.String() {
			// todo а оно работает?
			continue
		}

		if err = m.rep.Drop(id); err != nil {
			return false, err
		}
		delete(m.masks, id)
		//m.masks = append(m.masks[:i], m.masks[i+1:]...)
		return true, nil
	}

	return false, nil
}

func (m *MemoryIpGuard) Contains(ip string) (bool, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for id := range m.masks {
		if m.masks[id].Contains(net.ParseIP(ip)) {
			return true, nil
		}
	}
	return false, nil
}

func (m *MemoryIpGuard) Reload() error {
	var err error

	masks, err := m.rep.Get(m.listId)
	if err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()
	m.masks = make(map[int]*net.IPNet)
	for i := range masks {
		_, ipv4Net, err := net.ParseCIDR(masks[i].Mask)
		if err != nil {
			continue
		}
		m.masks[masks[i].Id] = ipv4Net
	}
	return err
}
