package model

import (
	"encoding/base64"
	"fmt"
	"testing"
)

const expected = `
                XXXXXXXXX                               XXXXXXXXX               
              XXXXXXXXXXXXX                           XXXXXXXXXXXXX             
             XXX         XXX                         XXX         XXX            
            XXX           XXX                       XXX           XXX           
           XX               XX                     XX               XX          
          XXX XXXXX         XXX                   XXX XXXXX         XXX         
          XX XXXXXXX         XX                   XX XXXXXXX         XX         
         XX  XXXXXXX          XX                 XX  XXXXXXX          XX        
         XX  XXXXXXX          XX                 XX  XXXXXXX          XX        
         XX  XXXXXXX          XX                 XX  XXXXXXX          XX        
         XX  XXXXXXX          XX                 XX  XXXXXXX          XX        
         XX   XXXXX           XX                 XX   XXXXX           XX        
         XX                   XX                 XX                   XX        
         XX                   XX                 XX                   XX        
         XX                   XX                 XX                   XX        
         XX                   XX                 XX                   XX        
          XX                 XX                   XX                 XX         
          XXX               XXX                   XXX               XXX         
           XX               XX                     XX               XX          
            XXX           XXX                       XXX           XXX           
             XXX         XXX                         XXX         XXX            
              XXXXXXXXXXXXX                           XXXXXXXXXXXXX             
                XXXXXXXXX                               XXXXXXXXX               
`

func TestModel(t *testing.T) {
	m := Model{
		mouseX:    10,
		mouseY:    5,
		winWidth:  80,
		winHeight: 25,
		debug:     false,
	}

	str := m.View()
	fmt.Println(str)

	if str != expected {
		t.Errorf("\nExpected:\n%s\nActual:\n%s", base64.StdEncoding.EncodeToString([]byte(expected)), base64.StdEncoding.EncodeToString([]byte(str)))
	}
}
