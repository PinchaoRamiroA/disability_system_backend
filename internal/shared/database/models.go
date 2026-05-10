package database

import (
	cobros "disability_system_backend/internal/modules/cobros/adapters/postgres/models"
	historial "disability_system_backend/internal/modules/historial/adapters/postgres/models"
	incapacidades "disability_system_backend/internal/modules/incapacidades/adapters/postgres/models"
	notificaciones "disability_system_backend/internal/modules/notificaciones/adapters/postgres/models"
	usuarios "disability_system_backend/internal/modules/usuarios/adapters/postgres/models"
)

var Models = []interface{}{
	// Usuarios
	&usuarios.UsuarioModel{},
	&usuarios.RolModel{},
	&usuarios.EmpleadoModel{},
	&usuarios.GerenciaModel{},
	&usuarios.GestionHumanaModel{},

	// Incapacidades
	&incapacidades.IncapacidadModel{},
	&incapacidades.DocumentoModel{},
	&incapacidades.EstadoIncapacidadModel{},
	&incapacidades.TipoIncapacidadModel{},
	&incapacidades.EntidadModel{},

	// Cobros
	&cobros.PagoModel{},
	&cobros.SeguimientoCobroModel{},
	&cobros.TipoSeguimientoModel{},

	// Historial
	&historial.HistorialModel{},
	&historial.TipoHistorialModel{},

	// Notificaciones
	&notificaciones.NotificacionModel{},
	&notificaciones.TipoNotificacionModel{},
}
