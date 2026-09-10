package temporal

import (
	v1 "api-registration-authorization/services/registration/contracts/v1"
	"api-registration-authorization/services/registration/dataaccess"
	"api-registration-authorization/services/registration/domain"
	"api-registration-authorization/shared"
	"context"
	"errors"

	"github.com/rs/zerolog/log"

	"github.com/google/uuid"
	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/client"
	sdktemporal "go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/worker"
	"google.golang.org/protobuf/proto"
)

type TemporalActivities struct {
	temporalClient client.Client
	worker         worker.Worker
	repository     *dataaccess.RegistrationRepository
}

func NewTemporalActivities(c client.Client, r *dataaccess.RegistrationRepository) *TemporalActivities {
	w := worker.New(c, domain.AuthQueue, worker.Options{})
	self := &TemporalActivities{
		temporalClient: c,
		worker:         w,
		repository:     r,
	}

	self.worker.RegisterActivityWithOptions(self.onCreateUser, activity.RegisterOptions{Name: domain.CreateUser})
	self.worker.RegisterActivityWithOptions(self.onDeleteUser, activity.RegisterOptions{Name: domain.DeleteUser})

	go func() {
		err := self.worker.Run(worker.InterruptCh())
		if err != nil {
			log.Fatal().Err(err).Msg("Temporal worker failed")
		}
	}()

	return self
}

func (self *TemporalActivities) onCreateUser(ctx context.Context, input *v1.CreateUserInput) (v1.CreateUserOutput, error) {
	activity.GetLogger(ctx).Info("User creating started....")

	userID, err := uuid.Parse(input.GetUserId())
	if err != nil {
		return v1.CreateUserOutput{
			Status:       v1.OperationStatus_OPERATION_STATUS_ERROR,
			ErrorMessage: proto.String("invalid auth user id"),
		}, nil
	}

	profileID, err := uuid.Parse(input.GetProfileId())
	if err != nil {
		return v1.CreateUserOutput{
			Status:       v1.OperationStatus_OPERATION_STATUS_ERROR,
			ErrorMessage: proto.String("invalid profile id"),
		}, nil
	}

	model := &shared.UserModel{
		Id:        userID,
		Email:     input.GetEmail(),
		Password:  input.GetHashedPassword(),
		ProfileId: profileID,
	}
	err = self.repository.Insert(model)

	switch {
	case err == nil:
		activity.GetLogger(ctx).Info("User creating done.")
		return v1.CreateUserOutput{
			Status: v1.OperationStatus_OPERATION_STATUS_SUCCESS,
		}, nil

	case errors.Is(err, domain.ErrorUserExists):
		return v1.CreateUserOutput{
			Status: v1.OperationStatus_OPERATION_STATUS_DUPLICATE,
		}, nil

	case errors.Is(err, domain.ErrorRegistrationConflict):
		return v1.CreateUserOutput{
			Status:       v1.OperationStatus_OPERATION_STATUS_ERROR,
			ErrorMessage: proto.String("registration id conflict"),
		}, nil

	default:
		return v1.CreateUserOutput{
			Status:      v1.OperationStatus_OPERATION_STATUS_SYSTEM_ERROR,
			SystemError: proto.Int64(-1),
		}, nil
	}
}

func (self *TemporalActivities) onDeleteUser(ctx context.Context, input *v1.DeleteUserInput) (v1.DeleteUserOutput, error) {
	activity.GetLogger(ctx).Info("User deleting started....")

	id, err := uuid.Parse(input.UserId)
	if err != nil {
		return v1.DeleteUserOutput{Success: false},
			sdktemporal.NewNonRetryableApplicationError(
				"invalid auth user id",
				"InvalidInput",
				err,
			)
	}

	err = self.repository.Delete(id)
	if err != nil {
		return v1.DeleteUserOutput{Success: false}, err
	}

	activity.GetLogger(ctx).Info("User deleting done.")
	return v1.DeleteUserOutput{Success: true}, nil
}
