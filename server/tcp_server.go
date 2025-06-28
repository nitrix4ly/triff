package server

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/nitrix4ly/triff/commands"
	"github.com/nitrix4ly/triff/core"
	"github.com/nitrix4ly/triff/utils"
)

// TCPServer handles TCP connections for Redis-like protocol
type TCPServer struct {
	db             *core.Database
	port           int
	listener       net.Listener
	stringCommands *commands.StringCommands
	logger         *utils.Logger
}

// NewTCPServer creates a new TCP server instance
func NewTCPServer(db *core.Database, port int, logger *utils.Logger) *TCPServer {
	return &TCPServer{
		db:             db,
		port:           port,
		stringCommands: commands.NewStringCommands(db),
		logger:         logger,
	}
}

// Start begins listening for TCP connections
func (s *TCPServer) Start() error {
	var err error
	s.listener, err = net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return fmt.Errorf("failed to start TCP server: %v", err)
	}

	s.logger.Info(fmt.Sprintf("TCP server listening on port %d", s.port))

	for {
		conn, err := s.listener.Accept()
		if err != nil {
			s.logger.Error(fmt.Sprintf("Error accepting connection: %v", err))
			continue
		}

		go s.handleConnection(conn)
	}
}

// Stop stops the TCP server
func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}

// handleConnection processes individual client connections
func (s *TCPServer) handleConnection(conn net.Conn) {
	defer conn.Close()
	
	s.logger.Info(fmt.Sprintf("New client connected: %s", conn.RemoteAddr()))
	
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		
		response := s.processCommand(line)
		conn.Write([]byte(response + "\r\n"))
	}
	
	if err := scanner.Err(); err != nil {
		s.logger.Error(fmt.Sprintf("Connection error: %v", err))
	}
	
	s.logger.Info(fmt.Sprintf("Client disconnected: %s", conn.RemoteAddr()))
}

// processCommand parses and executes commands
func (s *TCPServer) processCommand(input string) string {
	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "-ERR empty command"
	}
	
	command := strings.ToUpper(parts[0])
	args := parts[1:]
	
	switch command {
	case "PING":
		return "+PONG"
		
	case "SET":
		if len(args) < 2 {
			return "-ERR wrong number of arguments for 'set' command"
		}
		key, value := args[0], args[1]
		var ttl int64 = 0
		
		// Check for EX option (expiration in seconds)
		if len(args) >= 4 && strings.ToUpper(args[2]) == "EX" {
			var err error
			ttl, err = strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return "-ERR invalid expire time"
			}
		}
		
		response := s.stringCommands.Set(key, value, ttl)
		if response.Success {
			return "+OK"
		}
		return fmt.Sprintf("-ERR %s", response.Error)
		
	case "GET":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'get' command"
		}
		response := s.stringCommands.Get(args[0])
		if response.Success && response.Data != nil {
			return fmt.Sprintf("$%d\r\n%s", len(response.Data.(string)), response.Data.(string))
		}
		return "$-1"
		
	case "DEL":
		if len(args) == 0 {
			return "-ERR wrong number of arguments for 'del' command"
		}
		count := 0
		for _, key := range args {
			if s.db.Delete(key) {
				count++
			}
		}
		return fmt.Sprintf(":%d", count)
		
	case "EXISTS":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'exists' command"
		}
		if s.db.Exists(args[0]) {
			return ":1"
		}
		return ":0"
		
	case "KEYS":
		pattern := "*"
		if len(args) > 0 {
			pattern = args[0]
		}
		keys := s.db.Keys(pattern)
		result := fmt.Sprintf("*%d\r\n", len(keys))
		for _, key := range keys {
			result += fmt.Sprintf("$%d\r\n%s\r\n", len(key), key)
		}
		return result
		
	case "FLUSHALL":
		s.db.FlushAll()
		return "+OK"
		
	case "INFO":
		info := s.db.Info()
		result := ""
		for key, value := range info {
			result += fmt.Sprintf("%s:%v\r\n", key, value)
		}
		return fmt.Sprintf("$%d\r\n%s", len(result), result)
		
	case "DBSIZE":
		size := s.db.Size()
		return fmt.Sprintf(":%d", size)
		
	case "TTL":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'ttl' command"
		}
		ttl := s.db.GetTTL(args[0])
		return fmt.Sprintf(":%d", ttl)
		
	case "EXPIRE":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'expire' command"
		}
		seconds, err := strconv.ParseInt(args[1], 10, 64)
		if err != nil {
			return "-ERR invalid expire time"
		}
		if s.db.SetTTL(args[0], seconds) {
			return ":1"
		}
		return ":0"
		
	case "INCR":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'incr' command"
		}
		response := s.stringCommands.Incr(args[0])
		if response.Success {
			return fmt.Sprintf(":%d", response.Data.(int64))
		}
		return fmt.Sprintf("-ERR %s", response.Error)
		
	case "DECR":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'decr' command"
		}
		response := s.stringCommands.Decr(args[0])
		if response.Success {
			return fmt.Sprintf(":%d", response.Data.(int64))
		}
		return fmt.Sprintf("-ERR %s", response.Error)
		
	case "APPEND":
		if len(args) != 2 {
			return "-ERR wrong number of arguments for 'append' command"
		}
		response := s.stringCommands.Append(args[0], args[1])
		if response.Success {
			return fmt.Sprintf(":%d", response.Data.(int))
		}
		return fmt.Sprintf("-ERR %s", response.Error)
		
	case "STRLEN":
		if len(args) != 1 {
			return "-ERR wrong number of arguments for 'strlen' command"
		}
		response := s.stringCommands.Strlen(args[0])
		if response.Success {
			return fmt.Sprintf(":%d", response.Data.(int))
		}
		return fmt.Sprintf("-ERR %s", response.Error)
		
	default:
		return fmt.Sprintf("-ERR unknown command '%s'", command)
	}
}
