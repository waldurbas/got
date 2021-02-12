package dat

import (
	"fmt"
	"time"
)

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
// 2020.07.19 (wu) MkDate
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
		if schaltJahr(i) {
			days++
		}
	}

	for i = 1; i < mm; i++ {
		days += MonDays(i, yy)
	}

	days += dd

	return days
}

// PutDays #
func PutDays(days int) int {
	yy := 1980
	mm := 1
	dd := 1

	for days > 365 {
		tt := 365
		if schaltJahr(yy) {
			tt++
		}

		days -= tt

		if days == 0 {
			mm = 12
			dd = 31
		} else {
			yy++
		}
	}

	for days > 0 {
		for i := 1; days > 0 && i < 13; i++ {
			tt := MonDays(i, yy)
			if days > tt {
				mm = i + 1
				days -= tt
			} else if days > 0 {
				dd = days
				days = 0
			}
		}
	}

	return (yy * 10000) + (mm * 100) + dd
}

// MkDate #
func MkDate(dat int, dif int) int {
	return PutDays(GetDays(dat) + dif)
}

func schaltJahr(j int) bool {
	return (j%4 == 0 && j%100 != 0) || j%400 == 0
}

// PrintTimeDuration #
func PrintTimeDuration(dt time.Duration) string {
	inSek := int(dt / 1000000000)
	hh := inSek / 3600
	mm := inSek % 3600
	ss := mm % 60
	mm /= 60

	dd := hh / 24
	if dd > 0 {
		hh %= 24
	}

	stime := fmt.Sprintf("%.2d:%.2d:%.2d", hh, mm, ss)
	if dd > 0 {
		stime = fmt.Sprintf("%dd %.2d:%.2d:%.2d", dd, hh, mm, ss)
	}

	return stime
}
