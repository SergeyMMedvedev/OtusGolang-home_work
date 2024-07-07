package grpcserver

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ValidateYear(year int32) error {
	if year < 1970 || year > 9999 {
		return status.Errorf(codes.InvalidArgument, "Invalid year")
	}
	return nil
}

func ValidateMonth(month int32) error {
	if month < 1 || month > 12 {
		return status.Errorf(codes.InvalidArgument, "Invalid month")
	}
	return nil
}

func ValidateDay(day int32, month int32, year int32) error {
	if day < 1 || int(day) > time.Date(int(year), time.Month(month)+1, 0, 0, 0, 0, 0, time.UTC).Day() {
		return status.Errorf(codes.InvalidArgument, "Invalid day")
	}
	return nil
}
