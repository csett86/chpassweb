package main

import (
	"github.com/msteinert/pam"
)

// runPam starts a PAM chauthtok transaction in a goroutine and communicates via channels
func runPam(session *Session) {
	defer close(session.PromptCh)
	defer close(session.ResponseCh)

	// Conversation callback: pass PAM messages directly to the user
	conv := pam.ConversationFunc(func(style pam.Style, msg string) (string, error) {
		session.PromptCh <- msg
		return <-session.ResponseCh, nil
	})

	// Use "passwd" service as specified
	txn, err := pam.Start("passwd", session.Username, conv)
	if err != nil {
		session.ResultCh <- err
		return
	}

	err = txn.ChangeAuthTok(pam.ChangeExpiredAuthtok)
	session.ResultCh <- err
}
