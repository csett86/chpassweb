package main

import (
	"errors"
	"sync"
	"testing"
	"time"
)

func TestRunPam_Flow(t *testing.T) {
	// Skip this test as it requires PAM to be configured with root privileges
	t.Skip("Skipping: requires PAM configuration and root privileges")
	
	session := &Session{
		Username:   "testuser",
		PromptCh:   make(chan string, 10),
		ResponseCh: make(chan string),
		ResultCh:   make(chan error),
		CreatedAt:  time.Now(),
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		runPam(session)
	}()

	select {
	case prompt := <-session.PromptCh:
		t.Logf("Received PAM prompt: %s", prompt)
		session.ResponseCh <- "current_password"
		select {
		case <-session.PromptCh:
			session.ResponseCh <- "new_password"
			select {
			case <-session.PromptCh:
				session.ResponseCh <- "new_password"
				select {
				case <-session.ResultCh:
					t.Log("PAM finished")
				case <-time.After(100 * time.Millisecond):
					t.Fatal("Timeout waiting for PAM result")
				}
			case <-time.After(100 * time.Millisecond):
				t.Fatal("Timeout waiting for second PAM prompt")
			}
		case <-time.After(100 * time.Millisecond):
			t.Fatal("Timeout waiting for PAM response")
		}
	case <-time.After(100 * time.Millisecond):
		t.Fatal("Timeout waiting for first PAM prompt")
	}

	wg.Wait()
}

func MockRunPamWithError(session *Session, err error) {
	defer close(session.PromptCh)
	defer close(session.ResponseCh)
	session.ResultCh <- err
}

func TestErrorPropagation(t *testing.T) {
	session := &Session{
		Username:   "testuser",
		PromptCh:   make(chan string, 10),
		ResponseCh: make(chan string),
		ResultCh:   make(chan error),
		CreatedAt:  time.Now(),
	}

	mockError := errors.New("PAM authentication failed")
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		MockRunPamWithError(session, mockError)
	}()

	select {
	case err := <-session.ResultCh:
		if err.Error() != mockError.Error() {
			t.Errorf("Expected error '%v', got '%v'", mockError, err)
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Timeout waiting for error")
	}

	wg.Wait()
}
