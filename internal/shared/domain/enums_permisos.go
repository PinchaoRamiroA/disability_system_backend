package domain

type Permiso string

const (
	PermisoCrearIncapacidad     Permiso = "crear_incapacidad"
	PermisoEditarIncapacidad   Permiso = "editar_incapacidad"
	PermisoConsultarIncapacidad Permiso = "consultar_incapacidad"
	PermisoValidarDocumentos   Permiso = "validar_documentos"
	PermisoRechazarDocumentos   Permiso = "rechazar_documentos"
	PermisoRegistrarPago       Permiso = "registrar_pago"
	PermisoConsultarReportes   Permiso = "consultar_reportes"
	PermisoGestionarUsuarios   Permiso = "gestionar_usuarios"
	PermisoGestionarRoles      Permiso = "gestionar_roles"
	PermisoGenerarAlertas      Permiso = "generar_alertas"
	PermisoConsultarHistorial  Permiso = "consultar_historial"
	PermisoRealizarConciliacion Permiso = "realizar_conciliacion"
	PermisoCobroPersuasivo     Permiso = "gestionar_cobro_persuasivo"
	PermisoCobroJuridico       Permiso = "gestionar_cobro_juridico"
	PermisoArchivarIncapacidad  Permiso = "archivar_incapacidad"
)

func (p Permiso) IsValid() bool {
	permisos := []Permiso{
		PermisoCrearIncapacidad, PermisoEditarIncapacidad, PermisoConsultarIncapacidad,
		PermisoValidarDocumentos, PermisoRechazarDocumentos, PermisoRegistrarPago,
		PermisoConsultarReportes, PermisoGestionarUsuarios, PermisoGestionarRoles,
		PermisoGenerarAlertas, PermisoConsultarHistorial, PermisoRealizarConciliacion,
		PermisoCobroPersuasivo, PermisoCobroJuridico, PermisoArchivarIncapacidad,
	}
	for _, perm := range permisos {
		if p == perm {
			return true
		}
	}
	return false
}

func AllPermisos() []Permiso {
	return []Permiso{
		PermisoCrearIncapacidad, PermisoEditarIncapacidad, PermisoConsultarIncapacidad,
		PermisoValidarDocumentos, PermisoRechazarDocumentos, PermisoRegistrarPago,
		PermisoConsultarReportes, PermisoGestionarUsuarios, PermisoGestionarRoles,
		PermisoGenerarAlertas, PermisoConsultarHistorial, PermisoRealizarConciliacion,
		PermisoCobroPersuasivo, PermisoCobroJuridico, PermisoArchivarIncapacidad,
	}
}