package http

import (
	"context"
	"errors"

	"disability_system_backend/internal/modules/incapacidades/ports"
	apperrors "disability_system_backend/internal/shared/errors"

	"github.com/gin-gonic/gin"
)

type actorKey struct{}

func actorFromGin(c *gin.Context) (ports.Actor, error) {
	actorValue, exists := c.Get("actor")
	if !exists {
		return ports.Actor{}, apperrors.ErrUnauthorized.WithMessage("actor no encontrado en contexto")
	}
	actor, ok := actorValue.(ports.Actor)
	if !ok {
		return ports.Actor{}, apperrors.ErrInternal.WithMessage("actor inválido")
	}
	return actor, nil
}

func contextWithActor(ctx context.Context, actor ports.Actor) context.Context {
	return context.WithValue(ctx, actorKey{}, actor)
}

func actorFromContext(ctx context.Context) (ports.Actor, error) {
	actor, ok := ctx.Value(actorKey{}).(ports.Actor)
	if !ok {
		return ports.Actor{}, errors.New("actor not found")
	}
	return actor, nil
}
