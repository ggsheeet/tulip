// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package component

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Footer() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"footer_container\"><div class=\"footer_logo_tnc\"><div><a classhref=\"/\"><img src=\"https://storage.googleapis.com/tulip-storage/logos/logoreg-light-vert.webp\"></a></div><ul><li><a href=\"/\" class=\"nav_link\">Aviso de Privacidad</a></li><li><a href=\"/\" class=\"nav_link\">Copyright 2024</a></li><li><a href=\"/\" class=\"nav_link\">All Rights Reserved</a></li></ul></div><hr><div class=\"footer_routes_social\"><ul><li><a href=\"/\" class=\"nav_link\">Inicio</a></li><li><a href=\"/resources\" class=\"nav_link\">Recursos</a></li><li><a href=\"/store\" class=\"nav_link\">Tienda</a></li><li><a href=\"/articles\" class=\"nav_link\">Artículos</a></li><li><a href=\"/cart\" class=\"nav_link\">Carrito</a></li></ul><div class=\"footer_social\"><a href=\"https://www.facebook.com/profile.php?id=100076446663077\" rel=\"noopener noreferrer\" target=\"_blank\" class=\"social_img\"><img src=\"public/icons/facebook.svg\"></a> <a href=\"mailto:contacto@publicacionestulip.com\" class=\"social_img\"><img src=\"public/icons/email.svg\"></a></div></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
