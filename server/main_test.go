package main

import (
    "context"
    "testing"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"
    pb "golang-gRPC-user-service/proto/userpb"
    "github.com/stretchr/testify/assert"
)

func setupServer() *server {
    return &server{
        users: []pb.User{
            {Id: 1, Fname: "Ram", City: "Bhopal", Phone: 12345676, Height: 5.6, Married: true},
            {Id: 2, Fname: "Mayank", City: "Indore", Phone: 576768745, Height: 5.7, Married: true},
            {Id: 3, Fname: "Sheetal", City: "Pune", Phone: 44647687, Height: 5.4, Married: true},
            {Id: 4, Fname: "Riya", City: "Pune", Phone: 354645656, Height: 5.2, Married: true},
            {Id: 5, Fname: "Nainy", City: "Mumbai", Phone: 232432435, Height: 5.1, Married: true},
        },
    }
}

func TestGetUserDetails(t *testing.T) {
    srv := setupServer()

    tests := []struct {
        name      string
        request   *pb.UserIDRequest
        wantError bool
        wantUser  *pb.User
    }{
        {
            name:      "existing user",
            request:   &pb.UserIDRequest{Id: 1},
            wantError: false,
            wantUser:  &pb.User{Id: 1, Fname: "Ram", City: "Bhopal", Phone: 12345676, Height: 5.6, Married: true},
        },
        {
            name:      "non-existing user",
            request:   &pb.UserIDRequest{Id: 999},
            wantError: true,
            wantUser:  nil,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            res, err := srv.GetUserDetails(context.Background(), tt.request)
            if tt.wantError {
                assert.Error(t, err)
                assert.Nil(t, res)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, codes.NotFound, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, res)
                assert.Equal(t, tt.wantUser, res.User)
            }
        })
    }
}

func TestGetUsersDetails(t *testing.T) {
    srv := setupServer()

    tests := []struct {
        name      string
        request   *pb.UserIDsRequest
        wantError bool
        wantUsers []*pb.User
    }{
        {
            name: "existing users",
            request: &pb.UserIDsRequest{Ids: []int32{1, 2}},
            wantError: false,
            wantUsers: []*pb.User{
                {Id: 1, Fname: "Ram", City: "Bhopal", Phone: 12345676, Height: 5.6, Married: true},
                {Id: 2, Fname: "Mayank", City: "Indore", Phone: 576768745, Height: 5.7, Married: true},
            },
        },
        {
            name: "non-existing users",
            request: &pb.UserIDsRequest{Ids: []int32{999, 1000}},
            wantError: true,
            wantUsers: nil,
        },
        {
            name: "mix of existing and non-existing users",
            request: &pb.UserIDsRequest{Ids: []int32{1, 999}},
            wantError: false,
            wantUsers: []*pb.User{
                {Id: 1, Fname: "Ram", City: "Bhopal", Phone: 12345676, Height: 5.6, Married: true},
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            res, err := srv.GetUsersDetails(context.Background(), tt.request)
            if tt.wantError {
                assert.Error(t, err)
                assert.Nil(t, res)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, codes.NotFound, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, res)
                assert.Equal(t, tt.wantUsers, res.Users)
            }
        })
    }
}

func TestSearchUsers(t *testing.T) {
    srv := setupServer()

    tests := []struct {
        name       string
        request    *pb.SearchRequest
        wantError  bool
        wantResult []*pb.User
    }{
        {
            name: "search by existing city and non-existing phone number",
            request: &pb.SearchRequest{
                City:    "Pune",
                Phone: 358645656,
            },
            wantError: true,
            wantResult: nil,
        },
		{
            name: "search by existing city and existing phone number",
            request: &pb.SearchRequest{
                City:    "Mumbai",
                Phone: 232432435,
            },
            wantError: false,
            wantResult: []*pb.User{
                {Id: 5, Fname: "Nainy", City: "Mumbai", Phone: 232432435, Height: 5.1, Married: true},
            },
        },
        {
            name: "search by exising city",
            request: &pb.SearchRequest{
                City: "Indore",
            },
            wantError: false,
            wantResult: []*pb.User{
                {Id: 2, Fname: "Mayank", City: "Indore", Phone: 576768745, Height: 5.7, Married: true},
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            res, err := srv.SearchUsers(context.Background(), tt.request)
            if tt.wantError {
                assert.Error(t, err)
                assert.Nil(t, res)
                st, ok := status.FromError(err)
                assert.True(t, ok)
                assert.Equal(t, codes.NotFound, st.Code())
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, res)
                assert.Equal(t, tt.wantResult, res.Users)
            }
        })
    }
}
