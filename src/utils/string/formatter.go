package string

// import "strings"

// // According to the docs ()
// // The following codes are used to format the names and descriptions of servers
// // ^r 	reset
// // ^p 	newline (descriptions only)
// // ^n 	underline
// // ^l 	bold
// // ^m 	strike-through
// // ^o 	italic
// // ^0 	black
// // ^1 	blue
// // ^2 	green
// // ^3 	light blue
// // ^4 	red
// // ^5 	pink
// // ^6 	orange
// // ^7 	grey
// // ^8 	dark grey
// // ^9 	light purple
// // ^a 	light green
// // ^b 	light blue
// // ^c 	dark orange
// // ^d 	light pink
// // ^e 	yellow
// // ^f 	white

// type ColorMarker struct {
// 	color  uint8
// 	start  uint8
// 	end    uint8
// 	resets bool
// }

// type BeamFormatter struct {
// 	s string
// 	m []ColorMarker
// }

// // func FromTerminal(s string) BeamFormatter {

// // 	return BeamFormatter{s}
// // }

// func FromBeam(s string) BeamFormatter {
// 	fixed := ""
// 	colors := []ColorMarker{}
// 	for _, section := range strings.Split(s, "^") {
// 		if len(section) == 0 {
// 			fixed += "^"
// 			continue
// 		}
// 		switch section[0] {
// 		case 'r':
// 			colors = append(colors, ColorMarker{color: 0, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case 'p':
// 			colors = append(colors, ColorMarker{color: 0, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case 'n':
// 			colors = append(colors, ColorMarker{color: 1, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case 'l':
// 			colors = append(colors, ColorMarker{color: 2, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case 'm':
// 			colors = append(colors, ColorMarker{color: 3, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case 'o':
// 			colors = append(colors, ColorMarker{color: 4, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: false})
// 		case '0':
// 			colors = append(colors, ColorMarker{color: 5, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '1':
// 			colors = append(colors, ColorMarker{color: 6, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '2':
// 			colors = append(colors, ColorMarker{color: 7, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '3':
// 			colors = append(colors, ColorMarker{color: 8, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '4':
// 			colors = append(colors, ColorMarker{color: 9, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '5':
// 			colors = append(colors, ColorMarker{color: 10, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '6':
// 			colors = append(colors, ColorMarker{color: 11, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '7':
// 			colors = append(colors, ColorMarker{color: 12, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '8':
// 			colors = append(colors, ColorMarker{color: 13, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case '9':
// 			colors = append(colors, ColorMarker{color: 14, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'a':
// 			colors = append(colors, ColorMarker{color: 15, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'b':
// 			colors = append(colors, ColorMarker{color: 16, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'c':
// 			colors = append(colors, ColorMarker{color: 17, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'd':
// 			colors = append(colors, ColorMarker{color: 18, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'e':
// 			colors = append(colors, ColorMarker{color: 19, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		case 'f':
// 			colors = append(colors, ColorMarker{color: 20, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		default:
// 			colors = append(colors, ColorMarker{color: 0, start: uint8(len(fixed)), end: uint8(len(fixed)), resets: true})
// 		}
// 		fixed += section[1:]
// 	}
// 	return BeamFormatter{
// 		s: fixed,
// 		m: colors,
// 	}
// }

// func (f *BeamFormatter) ToTerminal() string {
// 	final := ""
// 	for _, c := range f.m {
// 		final += string(c.color)
// 	}
// 	return
// }

// func ParseForTerminal(s string) string {

// 	return s
// }
