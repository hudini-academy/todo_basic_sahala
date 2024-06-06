package models 
 
import ( 
    "errors" 
    "time" 
)

var ErrNoRecord = errors.New("models: no matching record found") 
 
type Todo struct { 
    ID      int 
    Title   string  
    Created time.Time 
    Expires time.Time 
}