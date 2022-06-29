package song

var TwelveBarBlues = Song{
	Title: "Twelve Bar Blues",
	Sections: Sections{
		ASection: [][]string{
			{"C", "C", "C", "C"},
			{"C", "C", "C", "C"},
			{"C", "C", "C", "C"},
			{"C", "C", "C", "C"},
		},
		BSection: [][]string{
			{"F", "F", "F", "F"},
			{"F", "F", "F", "F"},
			{"C", "C", "C", "C"},
			{"C", "C", "C", "C"},
		},
		CSection: [][]string{
			{"G", "G", "G", "G"},
			{"F", "F", "F", "F"},
			{"C", "C", "C", "C"},
			{"C", "C", "C", "C"},
		},
	},
}

var exptdTwlvBrBlsSctnA = []string{
	`C   
    
    
    `,
	` C  
    
    
    `,
	`  C 
    
    
    `,
	`   C
    
    
    `,
	`    
C   
    
    `,
	`    
 C  
    
    `,
	`    
  C 
    
    `,
	`    
   C
    
    `,
	`    
    
C   
    `,
	`    
    
 C  
    `,
	`    
    
  C 
    `,
	`    
    
   C
    `,
	`    
    
    
C   `,
	`    
    
    
 C  `,
	`    
    
    
  C `,
	`    
    
    
   C`,
}
var exptdTwlvBrBlsSctnB = []string{
	`F   
    
    
    `,
	` F  
    
    
    `,
	`  F 
    
    
    `,
	`   F
    
    
    `,
	`    
F   
    
    `,
	`    
 F  
    
    `,
	`    
  F 
    
    `,
	`    
   F
    
    `,
	`    
    
C   
    `,
	`    
    
 C  
    `,
	`    
    
  C 
    `,
	`    
    
   C
    `,
	`    
    
    
C   `,
	`    
    
    
 C  `,
	`    
    
    
  C `,
	`    
    
    
   C`,
}
var exptdTwlvBrBlsSctnC = []string{
	`G   
    
    
    `,
	` G  
    
    
    `,
	`  G 
    
    
    `,
	`   G
    
    
    `,
	`    
F   
    
    `,
	`    
 F  
    
    `,
	`    
  F 
    
    `,
	`    
   F
    
    `,
	`    
    
C   
    `,
	`    
    
 C  
    `,
	`    
    
  C 
    `,
	`    
    
   C
    `,
	`    
    
    
C   `,
	`    
    
    
 C  `,
	`    
    
    
  C `,
	`    
    
    
   C`,
}
var ExptectedTwelveBarBluesFrames = func() []string {
	ar := append(exptdTwlvBrBlsSctnA, exptdTwlvBrBlsSctnB...)
	arc := append(ar, exptdTwlvBrBlsSctnC...)
	return arc
}
