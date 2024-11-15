package server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/syedazeez337/note-taking-app-go/pb/github.com/syedazeez337/note-taking-app-go/pb"
	"google.golang.org/grpc"
)

// NoteServer implements pb.NoteServiceServer
type NoteServer struct {
	pb.UnimplementedNoteServiceServer
	notes map[string]*pb.Note
}

// NewNoteServer initializes a new NoteServer
func NewNoteServer() *NoteServer {
	return &NoteServer{
		notes: make(map[string]*pb.Note),
	}
}

// CreateNote creates a new note
func (s *NoteServer) CreateNote(ctx context.Context, req *pb.CreateNoteRequest) (*pb.CreateNoteResponse, error) {
	id := generateID() // Add a helper to generate unique IDs
	note := &pb.Note{
		Id:      id,
		Title:   req.Title,
		Content: req.Content,
	}
	s.notes[id] = note
	return &pb.CreateNoteResponse{Note: note}, nil
}

// GetNote retrieves a note by ID
func (s *NoteServer) GetNote(ctx context.Context, req *pb.GetNoteRequest) (*pb.GetNoteResponse, error) {
	note, exists := s.notes[req.Id]
	if !exists {
		return nil, errors.New("note not found")
	}
	return &pb.GetNoteResponse{Note: note}, nil
}

// ListNotes retrieves all notes
func (s *NoteServer) ListNotes(ctx context.Context, req *pb.ListNotesRequest) (*pb.ListNotesResponse, error) {
	notes := []*pb.Note{}
	for _, note := range s.notes {
		notes = append(notes, note)
	}
	return &pb.ListNotesResponse{Notes: notes}, nil
}

// UpdateNote updates an existing note
func (s *NoteServer) UpdateNote(ctx context.Context, req *pb.UpdateNoteRequest) (*pb.UpdateNoteResponse, error) {
	note, exists := s.notes[req.Id]
	if !exists {
		return nil, errors.New("note not found")
	}
	note.Title = req.Title
	note.Content = req.Content
	return &pb.UpdateNoteResponse{Note: note}, nil
}

// DeleteNote deletes a note by ID
func (s *NoteServer) DeleteNote(ctx context.Context, req *pb.DeleteNoteRequest) (*pb.DeleteNoteResponse, error) {
	_, exists := s.notes[req.Id]
	if !exists {
		return &pb.DeleteNoteResponse{Success: false}, errors.New("note not found")
	}
	delete(s.notes, req.Id)
	return &pb.DeleteNoteResponse{Success: true}, nil
}

// Helper function to generate unique IDs (use UUIDs in a real app)
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func StartServer() {
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	// create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the note service
	pb.RegisterNoteServiceServer(grpcServer, NewNoteServer())

	log.Println("Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server: %v:", err)
	}
}
