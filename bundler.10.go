//go:build generate

//go:generate fyne bundle --pkg iconres --prefix Icon -o ./ui/resources/iconres/signal_light_gray.gen.go ./ui/resources/iconres_src/signal_light_gray.png
//go:generate fyne bundle --pkg iconres --prefix Icon -o ./ui/resources/iconres/signal_light_green.gen.go ./ui/resources/iconres_src/signal_light_green.png
//go:generate fyne bundle --pkg iconres --prefix Icon -o ./ui/resources/iconres/signal_light_orange.gen.go ./ui/resources/iconres_src/signal_light_orange.png
//go:generate fyne bundle --pkg iconres --prefix Icon -o ./ui/resources/iconres/signal_light_red.gen.go ./ui/resources/iconres_src/signal_light_red.png
package solstice
