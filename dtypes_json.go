package main

import (
	"encoding/json"
	"errors"
)

type all struct {
	Type string `json:"type"`
}

type offerOrAnswer struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	SDP  string `json:"sdp"`
}

type register struct {
	Type   string `json:"type"`
	PeerID string `json:"peerID"`
}

type forward struct {
	Type string `json:"type"`
	From string `json:"from"`
	To   string `json:"to"`
	Data string `json:"data"`
}

func readOfferOrAnswer(msg []byte) (*offerOrAnswer, error) {
	var t offerOrAnswer
	err := json.Unmarshal(msg, &t)
	if err != nil {
		return nil, err
	}

	err = t.Check()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *offerOrAnswer) Check() error {
	if m.From == "" {
		return errors.New("no From field in JSON with type offerOrAnswer")
	}
	if m.To == "" {
		return errors.New("no To field in JSON with type offerOrAnswer")
	}
	if m.SDP == "" {
		return errors.New("no SDP field in JSON with type offerOrAnswer")
	}
	return nil
}

func readRegister(msg []byte) (*register, error) {
	var t register
	err := json.Unmarshal(msg, &t)
	if err != nil {
		return nil, err
	}

	err = t.Check()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *register) Check() error {
	if m.PeerID == "" {
		return errors.New("no PeerID field in JSON with type register")
	}
	return nil
}

func readForward(msg []byte) (*forward, error) {
	var t forward
	err := json.Unmarshal(msg, &t)
	if err != nil {
		return nil, err
	}

	err = t.Check()
	if err != nil {
		return nil, err
	}

	return &t, nil
}

func (m *forward) Check() error {
	if m.From == "" {
		return errors.New("no From field in JSON with type forward")
	}
	if m.To == "" {
		return errors.New("no To field in JSON with type forward")
	}
	if m.Data == "" {
		return errors.New("no Data field in JSON with type forward")
	}
	return nil
}
