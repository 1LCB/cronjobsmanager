# GoCronJobsManager

A simple and lightweight Go library for managing cron-like jobs using goroutines, `context.Context`, and a flexible scheduling system.

This library allows you to schedule jobs that:
- Run **once at a specific time**
- Run **repeatedly at a defined interval**
- **Start at a specific time** and then continue to run periodically

## Features

- Easy-to-use API for creating, listing, and removing jobs
- Context-based job cancellation
- Supports both one-time and recurring jobs
- Thread-safe job management with mutex locking
- Returns job ID for reference and control

---

## Installation

```bash
go get github.com/1lcb/cronjobsmanager
```

## Example

```golang
package main

import (
	"context"
	"fmt"
	"time"

	"github.com/1lcb/cronjobsmanager"
)

func main() {
	manager := cronjobsmanager.NewCronJobsManager()

	// Define a schedule: starts 5 seconds from now and repeats every 10 seconds
	schedule := cronjobsmanager.Schedule{
		StartAt: time.Now().Add(5 * time.Second),
		Every:   10 * time.Second,
	}

	// Create the job function
	jobFunc := func() {
		fmt.Println("Hello from scheduled job at", time.Now())
	}

	// Run the job
	jobId, err := manager.RunJob(context.Background(), "Say Hello", schedule, jobFunc)
	if err != nil {
		fmt.Println("Failed to start job:", err)
		return
	}

	fmt.Println("Job started with ID:", jobId)

	// Let it run for a while before stopping
	time.Sleep(35 * time.Second)

	// Remove the job
	err = manager.RemoveJob(jobId)
	if err != nil {
		fmt.Println("Failed to remove job:", err)
		return
	}

	fmt.Println("Job removed:", jobId)
}
```

