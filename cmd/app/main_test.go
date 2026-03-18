package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

const appAddr = "127.0.0.1:50051"

var (
	e2eEnabled    bool
	e2eSkipReason string

	appStartedByTests bool
	appCancel         context.CancelFunc
	appCmd            *exec.Cmd
	appLogs           bytes.Buffer
)

func TestMain(m *testing.M) {
	if waitForTCP(appAddr, 300*time.Millisecond) == nil {
		if err := preflightGRPC(appAddr); err == nil {
			e2eEnabled = true
			os.Exit(m.Run())
		}
		e2eSkipReason = "port 50051 is busy by non-compatible process"
		os.Exit(m.Run())
	}

	dbURL := strings.TrimSpace("postgres://postgres:test123@localhost:5432/postgres?sslmode=disable")
	if dbURL == "" {
		dbURL = strings.TrimSpace(os.Getenv("DATABASE_URL"))
	}
	if dbURL == "" {
		e2eSkipReason = "set E2E_DATABASE_URL or DATABASE_URL to run E2E tests"
		os.Exit(m.Run())
	}

	if err := startApp(dbURL); err != nil {
		e2eSkipReason = err.Error()
		os.Exit(m.Run())
	}

	e2eEnabled = true
	code := m.Run()
	stopApp()
	os.Exit(code)
}

func startApp(dbURL string) error {
	ctx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(ctx, "go", "run", ".")
	cmd.Dir = "."
	cmd.Env = append(os.Environ(), "DATABASE_URL="+dbURL)
	cmd.Stdout = &appLogs
	cmd.Stderr = &appLogs

	if err := cmd.Start(); err != nil {
		cancel()
		return fmt.Errorf("failed to start app: %w", err)
	}

	if err := waitForTCP(appAddr, 20*time.Second); err != nil {
		cancel()
		_ = cmd.Wait()
		return fmt.Errorf("app did not open %s: %w", appAddr, err)
	}

	if err := preflightGRPC(appAddr); err != nil {
		cancel()
		_ = cmd.Wait()
		return fmt.Errorf("grpc preflight failed: %w", err)
	}

	appStartedByTests = true
	appCancel = cancel
	appCmd = cmd
	return nil
}

func stopApp() {
	if !appStartedByTests || appCmd == nil {
		return
	}
	if appCancel != nil {
		appCancel()
	}

	done := make(chan error, 1)
	go func() {
		done <- appCmd.Wait()
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		if appCmd.Process != nil {
			_ = appCmd.Process.Kill()
			_ = appCmd.Wait()
		}
	}
}

func waitForTCP(addr string, timeout time.Duration) error {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		conn, err := net.DialTimeout("tcp", addr, 250*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			return nil
		}
		time.Sleep(100 * time.Millisecond)
	}
	return fmt.Errorf("timeout waiting for tcp %s", addr)
}

func preflightGRPC(addr string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return err
	}
	defer conn.Close()

	callCtx, callCancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer callCancel()

	err = conn.Invoke(callCtx, "/e2e.Preflight/UnknownMethod", &emptypb.Empty{}, &emptypb.Empty{})
	if err == nil {
		return fmt.Errorf("expected unimplemented error, got nil")
	}
	st, ok := status.FromError(err)
	if !ok || st.Code() != codes.Unimplemented {
		return fmt.Errorf("expected Unimplemented, got: %v", err)
	}
	return nil
}

func requireE2E(t *testing.T) {
	t.Helper()
	if !e2eEnabled {
		t.Skip(e2eSkipReason)
	}
}

func dialGRPC(t *testing.T) *grpc.ClientConn {
	t.Helper()

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, appAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		t.Fatalf("grpc dial failed: %v", err)
	}
	return conn
}

func TestE2E_ServerReachable(t *testing.T) {
	requireE2E(t)

	conn := dialGRPC(t)
	_ = conn.Close()
}

func TestE2E_ParallelTCPConnections(t *testing.T) {
	requireE2E(t)

	const n = 48
	for i := 0; i < n; i++ {
		i := i
		t.Run(fmt.Sprintf("tcp_%02d", i), func(t *testing.T) {
			t.Parallel()

			conn, err := net.DialTimeout("tcp", appAddr, 2*time.Second)
			if err != nil {
				t.Fatalf("tcp dial failed: %v", err)
			}
			_ = conn.Close()
		})
	}
}

func TestE2E_ParallelUnknownRPC(t *testing.T) {
	requireE2E(t)

	const n = 64
	for i := 0; i < n; i++ {
		i := i
		t.Run(fmt.Sprintf("rpc_%02d", i), func(t *testing.T) {
			t.Parallel()

			conn := dialGRPC(t)
			defer conn.Close()

			reqID := fmt.Sprintf("e2e-%d-%d", i, time.Now().UnixNano())
			method := fmt.Sprintf("/e2e.UnknownService/Method_%s", reqID)

			md := metadata.Pairs("x-e2e-id", reqID)
			ctx := metadata.NewOutgoingContext(context.Background(), md)

			callCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
			defer cancel()

			req := &wrapperspb.StringValue{Value: reqID}
			resp := &wrapperspb.StringValue{}

			err := conn.Invoke(callCtx, method, req, resp)
			if err == nil {
				t.Fatalf("expected error for unknown method, got nil")
			}
			st, ok := status.FromError(err)
			if !ok {
				t.Fatalf("non-status grpc error: %v", err)
			}
			if st.Code() != codes.Unimplemented {
				t.Fatalf("expected Unimplemented, got %s: %v", st.Code(), st.Message())
			}
		})
	}
}
