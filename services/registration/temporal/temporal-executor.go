package temporal

import (
	v1 "api-registration-authorization/services/registration/contracts/v1"
	"api-registration-authorization/services/registration/domain"
	"context"

	"github.com/rs/zerolog/log"
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

func (self *TemporalExecutor) BeginUserRegistration(email, hashedPassword string) (*v1.RegistrationResponse, error) {
	var result v1.RegistrationResponse

	log.Info().Str("workflow_type", domain.WorkFlowName).Msg("Begin user registration")

	input := &v1.RgistrationRequest{
		Email:        email,
		PasswordHash: hashedPassword,
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
		log.Error().Err(err).Msg("Registration workflow execution failed")
		return nil, err
	}

	err = run.Get(context.Background(), &result)
	if err != nil {
		log.Error().Err(err).Msg("Registration workflow result failed")
	}

	return &result, err
}
