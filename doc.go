/*
Package main - FlightGear MulitPlyer in go
	
	go-fgms
	------------------------------------------------------
	This is an experiment and silly idea of implementing the 
	FlightGear Multiplayer Server (fgms) in golang.

	There are a few challenges and the idea is to copy code line by line from c++ code to golang
	for fun.
	There are major difference and thats the fun ;-)!!! 
	"U can skin a cat in may ways"
	
	Currently is a step by step process, starting with
	- main
	- fg_config
	- fg_tracker
	- et all from there
	
	There are a few changes
	- isDeamon is gone, instead we expect this app to run with init.d supervisor, upstart or alike
	- tracker is in its own directory
	- simgear - might be a different project altogether
   
   
   Refs:
		- http://synflood.at/tmp/golang-slides/mrmcd2012.html# 
*/
package main