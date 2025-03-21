// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package component

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

import "github.com/ggsheet/tulip/internal/database"
import "strconv"
import "fmt"

func formatInts(id int) string {
	return strconv.Itoa(id)
}

func BookGrid(books []database.Book, currentPage int, totalPages int) templ.Component {
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
		if len(books) == 0 {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"books\" class=\"items_pagination\"><div class=\"items_container none\"><div class=\"none_message\" style=\"display: none;\"><p>No se encontraron Productos con este filtro</p><svg width=\"621\" height=\"530\" viewBox=\"0 0 621 530\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><style>\n                            .fill_bible_alt {\n                                fill: #484848;\n                            }\n                            @media (prefers-color-scheme: dark) {\n                                .fill_bible_alt {\n                                    fill: #afb6c0;\n                                }\n                            }\n                        </style><mask id=\"mask0_114_2\" style=\"mask-type:luminance\" maskUnits=\"userSpaceOnUse\" x=\"0\" y=\"70\" width=\"621\" height=\"461\"><path d=\"M0 70.8745H621V530H0V70.8745Z\" fill=\"white\"></path></mask> <g mask=\"url(#mask0_114_2)\"><path d=\"M620.501 82.6964V449.54H377.917L309.994 530.235L242.21 449.54H-0.263489V74.013H75.8456V373.561H277.72L309.994 411.931L342.532 373.561H544.406V74.013H620.501V82.6964ZM602.995 432.16V91.3935H561.788V390.942H350.469L309.994 439.239L269.52 390.942H58.215V91.3935H17.2427V432.035H250.534L310.119 502.927L369.842 432.035H602.995V432.16Z\" class=\"fill_bible_alt\"></path></g> <mask id=\"mask1_114_2\" style=\"mask-type:luminance\" maskUnits=\"userSpaceOnUse\" x=\"56\" y=\"0\" width=\"507\" height=\"381\"><path d=\"M56.6973 0H562.911V380.671H56.6973V0Z\" fill=\"white\"></path></mask> <g mask=\"url(#mask1_114_2)\"><path d=\"M58.8956 380.562V0.305969C132.142 0.305969 205.278 0.305969 278.401 0.305969L310.675 38.911L343.212 0.305969C416.335 0.305969 489.457 0.305969 562.593 0.305969V380.562H545.087V17.9354H351.15L310.675 65.9842L270.201 17.9354H76.5263V380.562H58.8956Z\" class=\"fill_bible_alt\"></path></g> <path d=\"M143.441 374.726L310.676 181.438L478.16 374.726L464.886 386.271L310.676 208.124L156.729 386.271L143.441 374.726Z\" class=\"fill_bible_alt\"></path> <path d=\"M61.1362 290.187L155.235 181.438L239.669 279.015L226.505 290.312L155.235 208.124L74.2866 301.733L61.1362 290.187Z\" class=\"fill_bible_alt\"></path> <path d=\"M541.727 301.733L460.902 208.124L389.647 290.312L376.358 279.015L460.902 181.438L554.891 290.187L541.727 301.733Z\" class=\"fill_bible_alt\"></path> <path d=\"M282.12 289.938H306.084V313.901H282.12V289.938ZM315.515 289.938H339.479V313.901H315.515V289.938ZM339.479 323.331V347.293H315.515V323.331H339.479ZM306.084 347.293H282.12V323.331H306.084V347.293Z\" class=\"fill_bible_alt\"></path> <path d=\"M301.991 86.0867H319.497V194.836H301.991V86.0867Z\" class=\"fill_bible_alt\"></path> <path d=\"M268.348 110.547H353.017V127.927H268.348V110.547Z\" class=\"fill_bible_alt\"></path></svg></div></div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		} else {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"books\" class=\"items_pagination\"><div class=\"items_container\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			for _, book := range books {
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"items_card\" style=\"display: none;\"><div class=\"card_cover\"><img src=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var2 string
				templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(book.CoverURL)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 56, Col: 52}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" alt=\"Book Cover\"></div><div class=\"card_info\"><div class=\"card_title\"><h2 class=\"title_clamp\">")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var3 string
				templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(book.Title)
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 60, Col: 68}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h2>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				if book.CoverType == "PDF" {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h4>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var4 string
					templ_7745c5c3_Var4, templ_7745c5c3_Err = templ.JoinStringErrs(book.CoverType)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 62, Col: 56}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var4))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h4>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				} else {
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<h4>Pasta ")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					var templ_7745c5c3_Var5 string
					templ_7745c5c3_Var5, templ_7745c5c3_Err = templ.JoinStringErrs(book.CoverType)
					if templ_7745c5c3_Err != nil {
						return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 64, Col: 62}
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var5))
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
					_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h4>")
					if templ_7745c5c3_Err != nil {
						return templ_7745c5c3_Err
					}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"card_price_qty\"><h5 data-price=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var6 string
				templ_7745c5c3_Var6, templ_7745c5c3_Err = templ.JoinStringErrs(formatPrice(book.Price))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 68, Col: 72}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var6))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">$")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var7 string
				templ_7745c5c3_Var7, templ_7745c5c3_Err = templ.JoinStringErrs(formatPrice(book.Price))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 68, Col: 101}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var7))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</h5><div class=\"card_qty\"><button class=\"decrement\" data-id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var8 string
				templ_7745c5c3_Var8, templ_7745c5c3_Err = templ.JoinStringErrs(formatInts(book.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 70, Col: 91}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var8))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" data-stock=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var9 string
				templ_7745c5c3_Var9, templ_7745c5c3_Err = templ.JoinStringErrs(formatInts(book.Stock))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 70, Col: 129}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var9))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><span>-</span></button><p id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var10 string
				templ_7745c5c3_Var10, templ_7745c5c3_Err = templ.JoinStringErrs("qty-" + formatInts(book.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 71, Col: 72}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var10))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">0</p><button class=\"increment\" data-id=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var11 string
				templ_7745c5c3_Var11, templ_7745c5c3_Err = templ.JoinStringErrs(formatInts(book.ID))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 72, Col: 91}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var11))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\" data-stock=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var12 string
				templ_7745c5c3_Var12, templ_7745c5c3_Err = templ.JoinStringErrs(formatInts(book.Stock))
				if templ_7745c5c3_Err != nil {
					return templ.Error{Err: templ_7745c5c3_Err, FileName: `template/component/bookgrid.templ`, Line: 72, Col: 129}
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var12))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\"><span>+</span></button></div></div><div class=\"card_view_add\"><a href=\"")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				var templ_7745c5c3_Var13 templ.SafeURL = templ.URL(fmt.Sprintf("/book?id=%d&category=%d", book.ID, book.CategoryID))
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(string(templ_7745c5c3_Var13)))
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
				_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("\">Ver más información</a> <button><svg width=\"256\" height=\"256\" viewBox=\"0 0 256 256\" fill=\"none\" xmlns=\"http://www.w3.org/2000/svg\"><style>\n                                            .fill_class {\n                                                fill: #f5f5f5;\n                                            }\n\n                                            .stroke_class {\n                                                stroke: #f5f5f5;\n                                            }\n\n                                            @media (prefers-color-scheme: dark) {\n                                                .fill_class {\n                                                    fill: #292929;\n                                                }\n\n                                                .stroke_class {\n                                                    stroke: #292929;\n                                                }\n                                            }\n                                        </style><path d=\"M53.3333 74.6666H138.667H200.427C213.048 74.6666 222.91 85.5641 221.654 98.1227L215.254 162.123C214.164 173.028 204.987 181.333 194.027 181.333H92.1558C81.9866 181.333 73.2311 174.156 71.2367 164.183L53.3333 74.6666Z\" stroke-width=\"25.6\" stroke-linejoin=\"round\" class=\"stroke_class\"></path> <path d=\"M53.3333 74.6667L44.6865 40.0796C43.4993 35.3312 39.2329 32 34.3383 32H21.3333\" stroke=\"#292929\" stroke-width=\"25.6\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"stroke_class\"></path> <path d=\"M85.3333 224H106.667\" stroke=\"#292929\" stroke-width=\"25.6\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"stroke_class\"></path> <path d=\"M170.667 224H192\" stroke=\"#292929\" stroke-width=\"25.6\" stroke-linecap=\"round\" stroke-linejoin=\"round\" class=\"stroke_class\"></path> <rect x=\"133.25\" y=\"93\" width=\"16.8269\" height=\"70\" rx=\"8.41346\" class=\"fill_class\"></rect> <rect x=\"107\" y=\"136.077\" width=\"16.8269\" height=\"70\" rx=\"8.41346\" transform=\"rotate(-90 107 136.077)\" class=\"fill_class\"></rect></svg></button></div></div></div>")
				if templ_7745c5c3_Err != nil {
					return templ_7745c5c3_Err
				}
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div><div class=\"pagination\">")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			templ_7745c5c3_Err = Pagination(currentPage, totalPages).Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div></div>")
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		return templ_7745c5c3_Err
	})
}
