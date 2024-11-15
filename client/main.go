package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/syedazeez337/note-taking-app-go/pb/github.com/syedazeez337/note-taking-app-go/pb"
	"google.golang.org/grpc"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewNoteServiceClient(conn)

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("\nAvailable Commands:")
		fmt.Println("1. create - Create a new note")
		fmt.Println("2. get - Get a note by ID")
		fmt.Println("3. list - List all notes")
		fmt.Println("4. update - Update a note")
		fmt.Println("5. delete - Delete a note by ID")
		fmt.Println("6. exit - Exit the client")

		fmt.Print("Enter command: ")
		command, _ := reader.ReadString('\n')
		command = strings.TrimSpace(command)

		switch command {
		case "create":
			createNote(client, reader)
		case "get":
			getNote(client, reader)
		case "list":
			listNotes(client)
		case "update":
			updateNote(client, reader)
		case "delete":
			deleteNote(client, reader)
		case "exit":
			fmt.Println("Exiting...")
			return
		default:
			fmt.Println("Invalid command, please try again.")
		}
	}
}

func createNote(client pb.NoteServiceClient, reader *bufio.Reader) {
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter content: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.CreateNote(ctx, &pb.CreateNoteRequest{
		Title:   title,
		Content: content,
	})
	if err != nil {
		log.Printf("Failed to create note: %v", err)
		return
	}

	fmt.Printf("Note created: %v\n", resp.Note)
}

func getNote(client pb.NoteServiceClient, reader *bufio.Reader) {
	fmt.Print("Enter note ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.GetNote(ctx, &pb.GetNoteRequest{Id: id})
	if err != nil {
		log.Printf("Failed to get note: %v", err)
		return
	}

	fmt.Printf("Note: %v\n", resp.Note)
}

func listNotes(client pb.NoteServiceClient) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.ListNotes(ctx, &pb.ListNotesRequest{})
	if err != nil {
		log.Printf("Failed to list notes: %v", err)
		return
	}

	fmt.Println("Notes:")
	for _, note := range resp.Notes {
		fmt.Printf("- %v\n", note)
	}
}

func updateNote(client pb.NoteServiceClient, reader *bufio.Reader) {
	fmt.Print("Enter note ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	fmt.Print("Enter new title: ")
	title, _ := reader.ReadString('\n')
	title = strings.TrimSpace(title)

	fmt.Print("Enter new content: ")
	content, _ := reader.ReadString('\n')
	content = strings.TrimSpace(content)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.UpdateNote(ctx, &pb.UpdateNoteRequest{
		Id:      id,
		Title:   title,
		Content: content,
	})
	if err != nil {
		log.Printf("Failed to update note: %v", err)
		return
	}

	fmt.Printf("Note updated: %v\n", resp.Note)
}

func deleteNote(client pb.NoteServiceClient, reader *bufio.Reader) {
	fmt.Print("Enter note ID: ")
	id, _ := reader.ReadString('\n')
	id = strings.TrimSpace(id)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	resp, err := client.DeleteNote(ctx, &pb.DeleteNoteRequest{Id: id})
	if err != nil {
		log.Printf("Failed to delete note: %v", err)
		return
	}

	if resp.Success {
		fmt.Println("Note deleted successfully.")
	} else {
		fmt.Println("Failed to delete note.")
	}
}
