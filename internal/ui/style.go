package ui

import (
	"fmt"
	"time"

	"github.com/julvo/htmlgo"
)

func FontScalingStyle(scale string) htmlgo.HTML {
	return htmlgo.Style_(htmlgo.Text(fmt.Sprintf(`
		html { 
			font-size: %s; 
			--text-scale: %s; /* CSS variable for consistency */
		}
		body { 
			font-size: inherit; 
		}
	`, scale, scale)))
}

func FormatDOB(dob string) string {
	dobTime, err := time.Parse("2006-01-02", dob)
	fmt.Printf("parsing DOB: dob=%s err=%v out=%v\n", dob, err, dobTime)
	if err == nil {
		return dobTime.Format("Jan 2, 2006")
	}
	return dob
}
