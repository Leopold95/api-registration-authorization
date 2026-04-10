package temporal

import (
	v1 "api-registration-authorization/services/registration/contracts/v1"
	"api-registration-authorization/services/registration/domain"
	"context"

	"go.temporal.io/sdk/client"
)

type TemporalExecutor struct {
	c client.Client
}

func NewTemporalExecutor(c client.Client) *TemporalExecutor {
	return &TemporalExecutor{
		c: c,
	}
}

func (self *TemporalExecutor) BeginUserRegistration(email, name, hashedPassword string) (*v1.RegistrationResponse, error) {
	var result v1.RegistrationResponse

	input := &v1.RgistrationRequest{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}

	run, err := self.c.ExecuteWorkflow(
		context.Background(),
		client.StartWorkflowOptions{
			ID:        "registration-" + email,
			TaskQueue: domain.OrchestratorQueue,
		},
		domain.WorkFlowName,
		input,
	)

	if err != nil {
		return nil, err
	}

	err = run.Get(context.Background(), &result)
	return &result, err
}
