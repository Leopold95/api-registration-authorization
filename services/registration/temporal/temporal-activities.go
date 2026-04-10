package temporal

import (
	v1 "api-registration-authorization/services/registration/contracts/v1"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/domain"
	"api-registration-authorization/shared"
	"context"
	"fmt"
	"log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

type TemporalActivities struct {
	temporalClient client.Client
	w              worker.Worker
	repository     *dataaccess.RegistrationRepository
}

func NewTemporalActivities(c client.Client, r *dataaccess.RegistrationRepository) *TemporalActivities {
	w := worker.New(c, domain.AuthQueue, worker.Options{})
	self := &TemporalActivities{
		temporalClient: c,
		w:              w,
		repository:     r,
	}

	self.w.RegisterActivityWithOptions(self.onCreateUser, activity.RegisterOptions{Name: domain.CreateUser})
	self.w.RegisterActivityWithOptions(self.onDeleteUser, activity.RegisterOptions{Name: domain.DeleteUser})

	go func() {
		err := self.w.Run(worker.InterruptCh())
		if err != nil {
			log.Fatal("cant start worker", err)
		}
	}()

	return self
}

func (self *TemporalActivities) onCreateUser(ctx context.Context, input *v1.CreateUserInput) (v1.CreateUserOutput, error) {
	fmt.Println("User creating started....")

	id, _ := uuid.Parse(input.ProfileId)
	model := &shared.UserModel{
		Email:     input.Email,
		Name:      input.Name,
		Password:  input.HashedPassword,
		ProfileId: id,
	}
	err := self.repository.Insert(model)

	if err != nil {
		return v1.CreateUserOutput{Success: false}, err
	}

	fmt.Println("User creating done.")
	return v1.CreateUserOutput{Success: true}, nil
}

func (self *TemporalActivities) onDeleteUser(ctx context.Context, input *v1.DeleteUserInput) (v1.DeleteUserOutput, error) {
	fmt.Println("User deleting started....")

	id, err := uuid.Parse(input.UserId)
	if err != nil {
		return v1.DeleteUserOutput{Success: false}, err
	}

	err = self.repository.Delete(id)
	if err != nil {
		return v1.DeleteUserOutput{Success: false}, err
	}

	fmt.Println("User deleting done.")
	return v1.DeleteUserOutput{Success: true}, nil
}
