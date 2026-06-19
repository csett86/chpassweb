package main

import (
	"github.com/msteinert/pam"
)

// runPam starts a PAM chauthtok transaction in a goroutine and communicates via channels
func runPam(session *Session) {
	defer close(session.PromptCh)
	defer close(session.ResponseCh)

	// Conversation callback: pass PAM messages directly to the user
	conv := pam.ConversationFunc(func(messages []pam.Message, numMsg int) ([]pam.Response, error) {
		for _, msg := range messages {
			session.PromptCh <- msg.Msg
		}
		response := <-session.ResponseCh
		return []pam.Response{{Resp: response}}, nil
	})

	// Use "passwd" service as specified
	txn, err := pam.Start("passwd", session.Username, conv)
	if err != nil {
		session.ResultCh <- err
		return
	}

	err = txn.Chauthtok(pam.ChangeExpireAuthTok)
	session.ResultCh <- err
}
