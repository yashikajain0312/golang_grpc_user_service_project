package main

import (
  "context"
  "log"
  "net"
  "google.golang.org/grpc"
  "google.golang.org/grpc/codes"
  "google.golang.org/grpc/status"
  pb "golang-gRPC-user-service/proto/userpb"
)

type server struct {
  pb.UnimplementedUserServiceServer
  users []pb.User
}

func (s *server) GetUserDetails(ctx context.Context, req *pb.UserIDRequest) (*pb.UserResponse, error) {
  for _, user := range s.users {
    if user.Id == req.Id {
      return &pb.UserResponse{User: &user}, nil
    }
  }
  return nil, status.Errorf(codes.NotFound, "User not found")
}

func (s *server) GetUsersDetails(ctx context.Context, req *pb.UserIDsRequest) (*pb.UsersResponse, error) {
  var foundUsers []*pb.User
  for _, id := range req.Ids {
    for _, user := range s.users {
      if user.Id == id {
        foundUsers = append(foundUsers, &user)
      }
    }
  }

  if len(foundUsers) == 0 {
	return nil, status.Errorf(codes.NotFound, "No user found with the given ids")
  }  

	return &pb.UsersResponse{Users: foundUsers}, nil
  
}

func (s *server) SearchUsers(ctx context.Context, req *pb.SearchRequest) (*pb.UsersResponse, error) {
  var foundUsers []*pb.User
  for _, user := range s.users {
    if (req.City == "" || user.City == req.City) &&
      (req.Phone == 0 || user.Phone == req.Phone) &&
      (!req.Married || user.Married == req.Married){
      foundUsers = append(foundUsers, &user)
    }
  }
  if len(foundUsers) == 0 {
	return nil, status.Errorf(codes.NotFound, "No user found")
}
	return &pb.UsersResponse{Users: foundUsers}, nil
}

func main() {
  lis, err := net.Listen("tcp", ":50051")
  if err != nil {
    log.Fatalf("failed to listen: %v", err)
  }

  s := grpc.NewServer()
  pb.RegisterUserServiceServer(s, &server{
    users: []pb.User{
      {Id: 1, Fname: "Ram", City: "Bhopal", Phone: 12345676, Height: 5.6, Married: true},
	  {Id: 2, Fname: "Mayank", City: "Indore", Phone: 576768745, Height: 5.7, Married: true},
	  {Id: 3, Fname: "Sheetal", City: "Pune", Phone: 44647687, Height: 5.4, Married: true},
	  {Id: 4, Fname: "Riya", City: "Pune", Phone: 354645656, Height: 5.2, Married: true},
	  {Id: 5, Fname: "Nainy", City: "Mumbai", Phone: 232432435, Height: 5.1, Married: true},
    },
  })

  log.Printf("server listening at %v", lis.Addr())
  if err := s.Serve(lis); err != nil {
    log.Fatalf("failed to serve: %v", err)
  }
}
