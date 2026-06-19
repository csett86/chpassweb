package main

import (
	"sync"
	"testing"
	"time"
)

func TestNewSessionStore(t *testing.T) {
	store := NewSessionStore(1 * time.Minute)
	if store == nil {
		t.Fatal("NewSessionStore returned nil")
	}
	if store.Timeout != 1*time.Minute {
		t.Errorf("Expected timeout 1m, got %v", store.Timeout)
	}
	if store.Sessions == nil {
		t.Fatal("Sessions map is nil")
	}
}

func TestSessionStore_Create(t *testing.T) {
	store := NewSessionStore(1 * time.Minute)
	session := store.Create("testuser")
	if session == nil {
		t.Fatal("Create returned nil session")
	}
	if session.Username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", session.Username)
	}
	if session.ID == "" {
		t.Error("Session ID is empty")
	}
	storedSession := store.Get(session.ID)
	if storedSession == nil {
		t.Error("Session not found in store")
	}
}

func TestSessionStore_Get(t *testing.T) {
	store := NewSessionStore(1 * time.Minute)
	if store.Get("nonexistent") != nil {
		t.Error("Expected nil for non-existent session")
	}
	created := store.Create("testuser")
	if store.Get(created.ID) == nil {
		t.Error("Expected to get created session")
	}
}

func TestSessionStore_Delete(t *testing.T) {
	store := NewSessionStore(1 * time.Minute)
	session := store.Create("testuser")
	if store.Get(session.ID) == nil {
		t.Error("Session should exist before delete")
	}
	store.Delete(session.ID)
	if store.Get(session.ID) != nil {
		t.Error("Session should not exist after delete")
	}
	store.Delete("nonexistent")
}

func TestGenerateSessionID(t *testing.T) {
	id1 := generateSessionID()
	id2 := generateSessionID()
	if id1 == "" {
		t.Error("generateSessionID returned empty string")
	}
	if id1 == id2 {
		t.Error("generateSessionID returned duplicate IDs")
	}
	if len(id1) != 32 {
		t.Errorf("Expected 32 char ID, got %d", len(id1))
	}
}

func TestSessionStore_Concurrency(t *testing.T) {
	store := NewSessionStore(1 * time.Minute)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			session := store.Create("user" + string(rune(id)))
			store.Get(session.ID)
			store.Delete(session.ID)
		}(i)
	}
	wg.Wait()
	store.Mu.RLock()
	defer store.Mu.RUnlock()
	if len(store.Sessions) != 0 {
		t.Errorf("Expected 0 sessions, got %d", len(store.Sessions))
	}
}
