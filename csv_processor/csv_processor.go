package main

import (
	"context"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"../services"

	"google.golang.org/grpc"
)

func readRecords(r io.Reader, handler func(services.Record) error) error {
	csvReader := csv.NewReader(r)
	fmt.Println("start reading file")
	for {
		rawRecord, err := csvReader.Read()
		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		if len(rawRecord) != 4 {
			continue
		}

		record := services.Record{
			ID:           rawRecord[0],
			Name:         rawRecord[1],
			Email:        rawRecord[2],
			MobileNumber: rawRecord[3],
		}

		if err = handler(record); err != nil {
			return err
		}
	}

	return nil
}

// CsvProcessorService represents a service that able to parse csv-files
type CsvProcessorService struct {
	dbService services.DatabaseServiceClient
}

func (s *CsvProcessorService) processFile(path string) error {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	readRecords(f, func(record services.Record) error {
		fmt.Printf("Sending record %#v\n", record)
		_, err := s.dbService.AddRecord(context.Background(), &record)
		return err
	})
	return nil
}

// ProcessFile processes csv-file and sends parsed records to the database service
func (s *CsvProcessorService) ProcessFile(ctx context.Context, request *services.CsvProcessingRequest) (*services.CsvProcessingResult, error) {
	// TODO: Open url on the internet
	err := s.processFile(request.URL)
	if err != nil {
		return nil, err
	}
	return &services.CsvProcessingResult{Status: "OK"}, nil
}

func main() {
	csvPath := flag.String("path", "", "path to csv file to process")
	flag.Parse()

	// TODO: Services configuration
	grcpConn, err := grpc.Dial("127.0.0.1:8081", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cant connect to grpc")
	}
	defer grcpConn.Close()

	if *csvPath == "" {
		log.Fatal("no file to process")
	}

	// TODO: make CsvProcessorService a server
	csvProcessor := CsvProcessorService{services.NewDatabaseServiceClient(grcpConn)}
	if err = csvProcessor.processFile(*csvPath); err != nil {
		log.Fatal(err)
	}
}
