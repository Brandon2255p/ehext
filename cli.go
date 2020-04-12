package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Brandon2255p/ehext/horizgen"
)

func main() {
	aggregateName := ""
	searchLocation := ""
	writeEventAppliers := false
	writeCommandHandlers := false
	flag.BoolVar(&writeEventAppliers, "appliers", false, "Overwrites the event appliers file")
	flag.BoolVar(&writeCommandHandlers, "handlers", false, "Overwrites the command handler class")
	flag.StringVar(&aggregateName, "aggregate", "", "the name of the aggregate. ThingAggregate")
	flag.StringVar(&searchLocation, "search", "", "the place to look for events")
	flag.Parse()
	if searchLocation == "" || aggregateName == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}
	directory := path.Dir(searchLocation)
	log.Printf("Using output directory: %s", directory)
	// making it very simple. No need to close the file.
	content, err := ioutil.ReadFile(searchLocation)
	if err != nil {
		log.Fatal(err)
	}

	// Convert []byte to string and print to screen
	text := string(content)
	classes := horizgen.ExtractClasses(text)
	var commands []horizgen.EventData
	var events []horizgen.EventData
	for i, class := range classes {
		name := horizgen.ExtractName(class)
		comment := horizgen.ExtractComment(class)
		log.Printf("%d of %d : %s", i, len(classes), name)
		data := horizgen.EventData{Name: name, Description: comment}
		if strings.HasSuffix(name, "Event") {
			events = append(events, data)
		}
		if strings.HasSuffix(name, "Command") {
			commands = append(commands, data)
		}
	}
	generated := horizgen.GenerateRegisterAggregate(aggregateName)
	outputFile := path.Join(directory, "genregisteraggregate.go")
	horizgen.Write(outputFile, generated)
	if len(commands) > 0 {
		{
			generated := horizgen.GenerateCommand(aggregateName, commands...)
			outputFile := path.Join(directory, "gencommandtypes.go")
			horizgen.Write(outputFile, generated)
		}
		{
			registerCommands := horizgen.GenerateRegisterCommand(commands)
			outputFile := path.Join(directory, "genregistercommands.go")
			log.Printf("Writing %s", outputFile)
			horizgen.Write(outputFile, registerCommands)
		}
		{
			handleCommands := horizgen.GenerateHandleCommand(aggregateName, commands)
			outputFile := path.Join(directory, "genhandlecommand.go")
			log.Printf("Writing %s", outputFile)
			horizgen.Write(outputFile, handleCommands)
		}
		if writeCommandHandlers {
			commandHandlers := horizgen.GenerateCommandHandlers(aggregateName, commands)
			outputFile := path.Join(directory, "commandhandlers.go")
			log.Printf("Writing %s", outputFile)
			horizgen.Write(outputFile, commandHandlers)
		}
	}
	if len(events) > 0 {
		{
			generated := horizgen.GenerateEvent(events...)
			outputFile := path.Join(directory, "geneventtypes.go")
			horizgen.Write(outputFile, generated)
		}
		{
			applyEvent := horizgen.GenerateApplyEvent(aggregateName, events)
			outputFile := path.Join(directory, "genapplyevent.go")
			log.Printf("Writing %s", outputFile)
			horizgen.Write(outputFile, applyEvent)
		}
		if writeEventAppliers {
			applyEvent := horizgen.GenerateEventAppliers(aggregateName, events)
			outputFile := path.Join(directory, "eventappliers.go")
			log.Printf("Writing %s", outputFile)
			horizgen.Write(outputFile, applyEvent)
		}
	}
}
