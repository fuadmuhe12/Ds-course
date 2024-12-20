package main

import (
	"bytes"
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/rpc"
	"os"
	"strconv"
)

type AddPaperArgs struct {
	Author  string
	Title   string
	Format  string
	Content []byte
}

type AddPaperReply struct {
	PaperNumber int
	Message     string
}

type ListPapersReply struct {
	Papers []struct {
		PaperNumber int
		Author      string
		Title       string
	}
}

type GetPaperArgs struct {
	PaperNumber int
}

type GetPaperDetailsReply struct {
	Author string
	Title  string
}

type FetchPaperReply struct {
	Content []byte
}

func main() {
	// Register gob types to handle binary data
	gob.Register([]byte{})

	// Define and parse command-line arguments
	command := flag.String("command", "", "Command to execute: add, list, detail, fetch")
	server := flag.String("server", "localhost:1234", "Server address")
	flag.Parse()

	if *command == "" {
		log.Fatalf("Error: Command is required. Use add, list, detail, or fetch.")
	}

	// Connect to the RPC server
	client, err := rpc.Dial("tcp", *server)
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer client.Close()

	// Handle commands
	switch *command {
	case "add":
		if len(flag.Args()) < 3 {
			log.Fatalf("Usage: --command=add --server=<server> <author> <title> <file_path>")
		}
		author := flag.Args()[0]
		title := flag.Args()[1]
		filePath := flag.Args()[2]
		addPaper(client, author, title, filePath)

	case "list":
		listPapers(client)

	case "detail":
		if len(flag.Args()) < 1 {
			log.Fatalf("Usage: --command=detail --server=<server> <paper_number>")
		}
		paperNumber, err := strconv.Atoi(flag.Args()[0])
		if err != nil {
			log.Fatalf("Invalid paper number: %v", err)
		}
		getPaperDetails(client, paperNumber)

	case "fetch":
		if len(flag.Args()) < 1 {
			log.Fatalf("Usage: --command=fetch --server=<server> <paper_number>")
		}
		paperNumber, err := strconv.Atoi(flag.Args()[0])
		if err != nil {
			log.Fatalf("Invalid paper number: %v", err)
		}
		fetchPaperContent(client, paperNumber)

	default:
		log.Fatalf("Unknown command: %s", *command)
	}
}

func addPaper(client *rpc.Client, author, title, filePath string) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read file: %v", err)
	}

	args := AddPaperArgs{
		Author:  author,
		Title:   title,
		Format:  detectFormat(filePath),
		Content: content,
	}
	var reply AddPaperReply

	err = client.Call("PaperServer.AddPaper", args, &reply)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Printf("Paper added successfully with PaperNumber: %d\n", reply.PaperNumber)
}

func listPapers(client *rpc.Client) {
	var reply ListPapersReply
	err := client.Call("PaperServer.ListPapers", struct{}{}, &reply)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Println("Stored Papers:")
	for _, paper := range reply.Papers {
		fmt.Printf("PaperNumber: %d, Author: %s, Title: %s\n", paper.PaperNumber, paper.Author, paper.Title)
	}
}

func getPaperDetails(client *rpc.Client, paperNumber int) {
	args := GetPaperArgs{PaperNumber: paperNumber}
	var reply GetPaperDetailsReply

	err := client.Call("PaperServer.GetPaperDetails", args, &reply)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Printf("Paper Details:\nAuthor: %s\nTitle: %s\n", reply.Author, reply.Title)
}

func fetchPaperContent(client *rpc.Client, paperNumber int) {
	args := GetPaperArgs{PaperNumber: paperNumber}
	var reply FetchPaperReply

	err := client.Call("PaperServer.FetchPaperContent", args, &reply)
	if err != nil {
		log.Fatalf("RPC error: %v", err)
	}

	fmt.Println("Paper Content:")
	fmt.Println(string(reply.Content))
}

func detectFormat(filePath string) string {
	switch {
	case bytes.HasSuffix([]byte(filePath), []byte(".pdf")):
		return "PDF"
	case bytes.HasSuffix([]byte(filePath), []byte(".doc")) || bytes.HasSuffix([]byte(filePath), []byte(".docx")):
		return "DOC"
	default:
		log.Fatalf("Unsupported file format. Use PDF or DOC.")
	}
	return ""
}
