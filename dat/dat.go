package dat

// ----------------------------------------------------------------------------------
// dat.go (https://github.com/waldurbas/got)
// Copyright 2020 by Waldemar Urbas
//-----------------------------------------------------------------------------------
// This Source Code Form is subject to the terms of the 'MIT License'
// A short and simple permissive license with conditions only requiring
// preservation of copyright and license notices.  Licensed works, modifications,
// and larger works may be distributed under different terms and without source code.
// ----------------------------------------------------------------------------------
// HISTORY
//-----------------------------------------------------------------------------------
// 2020.07.05 (wu) Init
//-----------------------------------------------------------------------------------

// MonDays # liefert Anzahl Tage im Monat/Jahr
func MonDays(mo int, ja int) int {
	switch mo {
	case 2:
		da := 28
		if (ja%4 == 0 && ja%100 != 0) || ja%400 == 0 {
			da++
		}

		return da
	case 4, 6, 9, 11:
		return 30
	}

	return 31
}

// GetDays # liefert Anzahl Tage seit 1980
func GetDays(dat int) int {
	i := dat % 10000

	yy := dat / 10000
	mm := i / 100
	dd := i % 100

	days := 0
	for i = 1980; i < yy; i++ {
		days += 365
		if (i%4 == 0 && i%100 != 0) || i%400 == 0 {
			days++
		}
	}

	for i = 1; i < mm; i++ {
		days += MonDays(i, yy)
	}

	days += dd

	return days
}
