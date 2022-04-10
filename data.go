package main

import (
	"time"

	"github.com/blackdatura/grpc.server/pb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var employees = []*pb.Employee{
	{
		Id:        0,
		No:        1994,
		FirstName: "Hill",
		LastName:  "Liu",
		MonthSalary: &pb.MonthSalary{
			Basic: 30000,
			Bonus: 15000,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
	{
		Id:        2,
		No:        1996,
		FirstName: "Jackey",
		LastName:  "Li",
		MonthSalary: &pb.MonthSalary{
			Basic: 20000,
			Bonus: 8000,
		},
		Status: pb.EmployeeStatus_NORMAL,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	}, {
		Id:        1,
		No:        1993,
		FirstName: "James",
		LastName:  "Wang",
		MonthSalary: &pb.MonthSalary{
			Basic: 18000,
			Bonus: 12000,
		},
		Status: pb.EmployeeStatus_ONVACATION,
		LastModified: &timestamppb.Timestamp{
			Seconds: time.Now().Unix(),
		},
	},
}
