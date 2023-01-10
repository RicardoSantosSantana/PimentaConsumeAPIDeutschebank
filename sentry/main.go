package main

import (
	"log"
	"time" 	
 
	"github.com/getsentry/sentry-go"
	structs "github.com/fatih/structs"
	   
)
type Tokenb struct {
	AccessToken  string 
	TokenType    string 
	RefreshToken string 
	ExpiresIn    int    
	WhenExpires  string 
	Scope         string 
	Name string
	Email string
}


func main() {

 
	err := sentry.Init(sentry.ClientOptions{
		//Dsn: "https:04c1ffa5ffbf4eebb3f392b6069d38aa@o4504417796882432.ingest.sentry.io/4504418327986176",
		Dsn: "http://ba3a0c58015c4b2da6ca16af0fed86b0@pimenta:9000/2",
		Environment: "Ellie",
		Release: "pimenta-finance@1.0.0",

	})

	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}
 
	defer sentry.Flush(2 * time.Second)

	token := &Tokenb{
		AccessToken:"04c1ffa5ffbf4eebb3f392b6069d38aa@o4504417796882432",
		TokenType:"Code",
		RefreshToken:"04c1ffa5ffbf4eebb3f392b6069d38aa@o4504417796882432Refresh",
		ExpiresIn:10,
		WhenExpires:"2022-12-30 19:30",
		Scope:"global",
		Name:"Ricardo Souza",
		Email:"r.souza@pimenta.group",
	} 

	usuario:=sentry.User{
		Username:"passarinho",
		Name:"Ricardo Souza",
		Email: "r.souza@gmail.com",
		Segment:"Devops",		
	}

	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetContext("Devops Requirements",structs.Map(token))
		scope.SetUser(usuario)
	})
   
/*
	const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warning"
	LevelError   Level = "error"
	LevelFatal   Level = "fatal"
)
	*/

	/*
	type Breadcrumb struct {
	Type      string                 `json:"type,omitempty"`
	Category  string                 `json:"category,omitempty"`
	Message   string                 `json:"message,omitempty"`
	Data      map[string]interface{} `json:"data,omitempty"`
	Level     Level                  `json:"level,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}
	*/
 
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Category: "Messages",
		Message: "Authenticated user " + "r.santana@pimenta.group",
		Level: sentry.LevelFatal,
		Data: structs.Map(token),
	});
 
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("Devops", "Ricardo Souza");
		scope.SetLevel(sentry.LevelDebug)
	})

	sentry.CaptureMessage( "There is a problem on create database fination" )
}