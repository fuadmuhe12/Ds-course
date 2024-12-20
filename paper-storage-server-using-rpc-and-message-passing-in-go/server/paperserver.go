package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"net/rpc"
	"sync"

	"github.com/streadway/amqp"
)

// Paper structure to represent each paper type
type Paper struct {
	PaperNumber int
	Author      string
	Title       string
	Format      string
	Content     []byte
}

type PaperServer struct {
	papers []Paper
	mutex  sync.Mutex
	queue  string
}

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

// AddPaper handles storing a paper in memory and sending a RabbitMQ notification
func (ps *PaperServer) AddPaper(args AddPaperArgs, reply *AddPaperReply) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	paper := Paper{
		PaperNumber: len(ps.papers) + 1,
		Author:      args.Author,
		Title:       args.Title,
		Format:      args.Format,
		Content:     args.Content,
	}

	ps.papers = append(ps.papers, paper)
	reply.PaperNumber = paper.PaperNumber
	reply.Message = "Paper added successfully."

	// Publish notification to RabbitMQ
	message := fmt.Sprintf("New paper added: %s by %s", paper.Title, paper.Author)
	err := publishMessage(ps.queue, message)
	if err != nil {
		log.Printf("Failed to publish RabbitMQ message: %v", err)
	}

	return nil
}

// ListPapers returns a list of all stored papers
func (ps *PaperServer) ListPapers(args struct{}, reply *ListPapersReply) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	for _, paper := range ps.papers {
		reply.Papers = append(reply.Papers, struct {
			PaperNumber int
			Author      string
			Title       string
		}{
			PaperNumber: paper.PaperNumber,
			Author:      paper.Author,
			Title:       paper.Title,
		})
	}
	return nil
}

// GetPaperDetails returns the author and title of a specified paper
func (ps *PaperServer) GetPaperDetails(args GetPaperArgs, reply *GetPaperDetailsReply) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if args.PaperNumber <= 0 || args.PaperNumber > len(ps.papers) {
		return fmt.Errorf("Paper not found")
	}

	paper := ps.papers[args.PaperNumber-1]
	reply.Author = paper.Author
	reply.Title = paper.Title
	return nil
}

// FetchPaperContent retrieves the full content of a paper
func (ps *PaperServer) FetchPaperContent(args GetPaperArgs, reply *FetchPaperReply) error {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	if args.PaperNumber <= 0 || args.PaperNumber > len(ps.papers) {
		return fmt.Errorf("Paper not found")
	}

	paper := ps.papers[args.PaperNumber-1]
	reply.Content = paper.Content
	return nil
}

// publishMessage sends a message to a RabbitMQ queue
func publishMessage(queueName, message string) error {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		return err
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(
		queueName, // queue name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	return ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		},
	)
}

func main() {
	gob.Register([]byte{})

	server := new(PaperServer)
	server.queue = "papers"

	rpc.Register(server)

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	defer listener.Close()

	log.Println("PaperServer is running on port 1234...")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go rpc.ServeConn(conn)
	}
}
