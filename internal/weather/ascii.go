package weather

// WeatherAsciiArt maps weather condition IDs to ASCII art representations
var WeatherAsciiArt = map[int]string{
	// Clear sky
	800: `
    \   /
     .-.
  ― (   ) ―
     '-'
    /   \
`,
	// Few clouds
	801: `
   \  /
 _ /"".-.
   \_(   ).
   /(___(__) 
`,
	// Scattered clouds
	802: `
     .--.
  .-(    ).
 (___.__)__)
`,
	// Broken clouds
	803: `
     .--.
  .-(    ).
 (___.__)__)
     *   *
`,
	// Shower rain
	500: `
     .-.
    (   ).
   (___(__)
    ' ' ' '
   ' ' ' '
`,
	// Rain
	501: `
     .-.
    (   ).
   (___(__)
  ‚'‚'‚'‚'
 ‚'‚'‚'‚'
`,
	// Thunderstorm
	200: `
     .-.
    (   ).
   (___(__)
  ⚡''⚡''
 '⚡''⚡'
`,
	// Snow
	600: `
     .-.
    (   ).
   (___(__)
    *  *  *
   *  *  *
`,
	// Mist
	701: `
 _ - _ - _ -
  _ - _ - _
 _ - _ - _ -
`,
}

// GetWeatherAscii returns the ASCII art for a given weather condition ID
func GetWeatherAscii(conditionID int) string {
	if art, ok := WeatherAsciiArt[conditionID]; ok {
		return art
	}
	// Default ASCII art for unknown condition
	return `
   ?????
  ?     ?
 ?       ?
  ?     ?
   ?????

 Sorry. This ASCII art is not ready yet.
`
}
