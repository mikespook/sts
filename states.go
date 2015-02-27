package sts

import "github.com/mikespook/sts/model"

func (s *Sts) Sessions() *model.Sessions {
	return s.sessions
}

func (s *Sts) Agents() *model.Agents {
	return s.agents
}
